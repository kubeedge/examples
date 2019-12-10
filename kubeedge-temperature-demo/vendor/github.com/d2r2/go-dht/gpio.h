//--------------------------------------------------------------------------------------------------
//
// Copyright (c) 2015-2019 Denis Dyakov
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
// associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial
// portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
// BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
//--------------------------------------------------------------------------------------------------

#ifndef GO_DHT_H
#define GO_DHT_H

#ifndef _GNU_SOURCE
#define _GNU_SOURCE 1
#endif

#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <fcntl.h>
#include <sched.h>
#include <time.h>
#include <unistd.h>

// Add general macro block to detect OS
#ifdef _WIN32
   //define something for Windows (32-bit and 64-bit, this part is common)
   #ifdef _WIN64
      //define something for Windows (64-bit only)
   #else
      //define something for Windows (32-bit only)
   #endif
#elif __APPLE__
    #include "TargetConditionals.h"
    #if TARGET_IPHONE_SIMULATOR
         // iOS Simulator
    #elif TARGET_OS_IPHONE
        // iOS device
    #elif TARGET_OS_MAC
        // Other kinds of Mac OS
    #else
    #   error "Unknown Apple platform"
    #endif
#elif __linux__
    // linux
#elif __unix__ // all unices not caught above
    // Unix
#elif defined(_POSIX_VERSION)
    // POSIX
#else
#   error "Unknown compiler"
#endif

// GPIO direction: receive either output data to specific GPIO pin.
#define IN  0
#define OUT 1
 
// LOW correspong to low level of output signal, HIGH correspond to high level.
#define LOW  0
#define HIGH 1

// TRUE, FALSE values
#define FALSE 0
#define TRUE 1

// Keep pin no, file descriptors for data reading/writing
// and for specifying input/output mode.
typedef struct {
    int pin;
    // Keep file descriptors for "direction" and "value"
    // open during whole sensor session interraction,
    // because it save milliseconds critical for one-wire
    // DHTxx protocol
    int fd_direction;
    int fd_value;
} Pin;

// Struct to keep error info on function return.
typedef struct {
    char *message;
} Error;

// Create error with formated message.
// Skip error creation if no variable provided either error already exists.
static void create_error(Error **err, const char* format, ...) {
    if (err != NULL && *err == NULL) {
        *err = (Error *)malloc(sizeof(Error));
        va_list argptr;
        va_start(argptr, format);
        // get return value to suppress warning warn_unused_result
        int res = vasprintf(&(*err)->message, format, argptr);
        va_end(argptr);
    }
}
 
static void free_error(Error *err) {
    if (err != NULL) {
        free(err->message);
    }
    free(err);
}

// char* get_error_message(char const *msg) {
//     size_t needed = snprintf(NULL, 0, "%s: %s (%d)", msg, strerror(errno), errno) + 1;
//     char  *buffer = malloc(needed);
//     snprintf(buffer, needed, "%s: %s (%d)", msg, strerror(errno), errno);
//     return buffer;
// }

// Freeze thread for usec microseconds.
static int sleep_usec(int32_t usec) {
    struct timespec tim, tim2;
    // convert microseconds to seconds
    tim.tv_sec = usec / 1000000;
    // rest part of microseconds convert to nanoseconds
    tim.tv_nsec = (usec % 1000000) * 1000;
    return nanosleep(&tim , &tim2);
}

// Start working with specific pin.
static int gpio_export(int port, Pin *pin, Error **err) {
    #define BUFFER_MAX 3
    char buffer[BUFFER_MAX];
    ssize_t bytes_written;
    int fd;

    // Initialize pin to work with, "direction" and "value"
    // file descriptors with empty value.
    pin->pin = -1;
    pin->fd_direction = -1;
    pin->fd_value = -1;
                 
    fd = open("/sys/class/gpio/export", O_WRONLY|O_SYNC|O_RSYNC);
    if (-1 == fd) {
        create_error(err, "failed to open GPIO export for writing");
        return -1;
    }
    pin->pin = port;
    bytes_written = snprintf(buffer, BUFFER_MAX, "%d", pin->pin);
    if (-1 == write(fd, buffer, bytes_written)) {
        create_error(err, "failed to export pin %d", pin->pin);
        close(fd);
        return -1;
    }
    close(fd);

    // !!! Found in experimental way, that additional pause should exist
    // between export pin to work with and direction set up. Otherwise,
    // under the regular user mistake occures frequently !!!
    //
    // Sleep 150 milliseconds
    // sleep_usec(150*1000);

    #define DIRECTION_MAX 35
    char path1[DIRECTION_MAX];
    snprintf(path1, DIRECTION_MAX, "/sys/class/gpio/gpio%d/direction", pin->pin);
    pin->fd_direction = open(path1, O_WRONLY|O_SYNC|O_RSYNC);
    if (-1 == pin->fd_direction) {
        create_error(err, "failed to open pin %d direction for writing", pin->pin);
        return -1;
    }
                             
    #define VALUE_MAX 30
    char path2[VALUE_MAX];
    snprintf(path2, VALUE_MAX, "/sys/class/gpio/gpio%d/value", pin->pin);
    pin->fd_value = open(path2, O_RDWR|O_SYNC|O_RSYNC);
    if (-1 == pin->fd_value) {
        create_error(err, "failed to open pin %d value for reading", pin->pin);
        return -1;
    }
                             
    return 0;
}

// Stop working with specific pin.
static int gpio_unexport(Pin *pin, Error **err) {
    // Close "direction" file descriptor.
    if (-1 != pin->fd_direction) {
        close(pin->fd_direction);
        pin->fd_direction = -1;
    }
    // Close "value" file descriptor.
    if (-1 != pin->fd_value) {
        close(pin->fd_value);
        pin->fd_value = -1;
    }

    if (-1 != pin->pin) {
        char buffer[BUFFER_MAX];
        ssize_t bytes_written;
        int fd;
                 
        fd = open("/sys/class/gpio/unexport", O_WRONLY|O_SYNC|O_RSYNC);
        if (-1 == fd) {
            create_error(err, "failed to open unexport for writing");
            return -1;
        }
                         
        bytes_written = snprintf(buffer, BUFFER_MAX, "%d", pin->pin);
        if (-1 == write(fd, buffer, bytes_written)) {
            create_error(err, "failed to unexport pin %d", pin->pin);
            close(fd);
            return -1;
        }

        close(fd);
    }
    return 0;
}
 
// Setup pin mode: input or output.
static int gpio_direction(Pin *pin, int dir, Error **err) {
    static const char s_directions_str[]  = "in\0out";
    if (-1 == write(pin->fd_direction, &s_directions_str[IN == dir ? 0:3],
            IN == dir ? 2:3)) {
        create_error(err, "failed to set direction \"%s\" to pin %d",
            &s_directions_str[IN == dir ? 0:3], pin->pin);
        return -1;
    }
    return 0;
}

// Read data from the pin: in normal conditions return 0 or 1,
// which correspond to low or high signal levels.
// Experimantally found, that data might be preterminated
// with line end ('\n').
static int gpio_read(Pin *pin, Error **err) {
    char value_str[3];
    // Seek and read in one call; use instead sequential lseek() and read().
    if (-1 == pread(pin->fd_value, value_str, 3, 0)) {
        create_error(err, "failed to read value");
        return -1;
    }

    // printf("bytes read: %02X:%02X:%02X\n", value_str[0], value_str[1], value_str[2]);

    // Small optimization to speed up GPIO processing to skip
    // atoi() due to ARM devices poor CPU performance.
    if ((value_str[1] == '\0' || value_str[1] == '\n') &&
           (value_str[0] == '0' || value_str[1] == '1')) {
        return value_str[0] == '0' ? 0:1;
    } else {
        char *pos = strchr(value_str, '\n');
        if (pos != NULL)
            *pos = '\0';
        return atoi(value_str);
    }
}

// Set up specific pin level to 0 (low) or 1 (high).
static int gpio_write(Pin *pin, int value, Error **err) {
    static const char s_values_str[] = {'0', 0, '1'};
    if (1 != write(pin->fd_value, &s_values_str[LOW == value ? 0:2], 1)) {
        create_error(err, "failed to write value \"%s\" to pin %d",
            &s_values_str[LOW == value ? 0:2], pin->pin);
        return -1;
    }
    return 0;
}

// Macro to convert timespec structure value to microseconds.
#define convert_timespec_to_usec(t) ((t).tv_sec*1000*1000 + (t).tv_nsec/1000)
 
// Read sequence of data from the pin trigering
// on edge change until timeout occures.
// Collect as well durations of pulses in microseconds.
// Fill [arr] array with a sequence: level1, duration1, level2, duration2...
// Put array length to variable [len].
static int gpio_read_seq_until_timeout(Pin *pin,
        int32_t timeout_msec, int32_t **arr, int32_t *len, Error **err) {
    int32_t last_v, next_v;
#define MAX_PULSE_COUNT 16000
    int values[MAX_PULSE_COUNT*2];
    
    last_v = gpio_read(pin, err);
    if (-1 == last_v) {
        create_error(err, "failed to read value");
        return -1;
    }
    int k = 0, i = 0;
    values[k*2] = last_v;
    struct timespec last_t, next_t;
#define CLOCK_KIND CLOCK_MONOTONIC
// #define CLOCK_KIND CLOCK_REALTIME
    clock_gettime(CLOCK_KIND, &last_t);

    for (;;)
    {
        next_v = gpio_read(pin, err);
        if (-1 == next_v) {
            create_error(err, "failed to read value");
            return -1;
        }

        // Check for edge trigger event.
        if (last_v != next_v) {
            clock_gettime(CLOCK_KIND, &next_t); 
            i = 0;
            k++;
            if (k > MAX_PULSE_COUNT-1) {
                create_error(err, "pulse count exceed limit in %d", MAX_PULSE_COUNT);
                return -1;
            }
            values[k*2] = next_v;
            // Save time duration in microseconds of last edge level.
            values[k*2-1] = convert_timespec_to_usec(next_t) -
                convert_timespec_to_usec(last_t); 
            last_v = next_v;
            last_t = next_t;
        }

        // Each N cycle, without edge trigger event,
        // try to detect expiration of timeout, to terminate processing.
        if (i++ > 30) {
            clock_gettime(CLOCK_KIND, &next_t);
            int dur = convert_timespec_to_usec(next_t) -
                convert_timespec_to_usec(last_t);
            if (dur / 1000 > timeout_msec) {
                values[k*2+1] = dur;
                break;
            }
        }
    }
    *arr = malloc((k+1)*2 * sizeof(int32_t));
    for (i=0; i<=k; i++)
    {
        (*arr)[i*2] = values[i*2];
        (*arr)[i*2+1] = values[i*2+1];
    }
    *len = (k+1)*2;
                                 
/*    fprintf(stdout, "scan %d values\n", k+1);
    for (i=0; i<=k; i++)
    {
        fprintf(stdout, "value %d (%d): %d\n", i, (*arr)[i*2+1], (*arr)[i*2]);
    }*/
    return 0;
}
 
 
#if !defined(__APPLE__) // sched_setscheduler() doesn't defined on Apple devices

// Used to gain maximum performance from device during
// receiving bunch of data from sensors like DHTxx.
static int set_max_priority(Error **err) {
    struct sched_param sched;
    memset(&sched, 0, sizeof(sched));
    // Use FIFO scheduler with highest priority
    // for the lowest chance of the kernel context switching.
    sched.sched_priority = sched_get_priority_max(SCHED_FIFO);
    if (-1 == sched_setscheduler(0, SCHED_FIFO, &sched)) {
        create_error(err, "unable to set SCHED_FIFO priority to the thread");
        return -1;
    }
    return 0;
}

// Get back normal thread priority.
static int set_default_priority(Error **err) {
    struct sched_param sched;
    memset(&sched, 0, sizeof(sched));
    // Go back to regular schedule priority.
    sched.sched_priority = 0;
    if (-1 == sched_setscheduler(0, SCHED_OTHER, &sched)) {
        create_error(err, "unable to set SCHED_OTHER priority to the thread");
        return -1;
    }
    return 0;
}

#endif // sched_setscheduler() doesn't defined on Apple devices


// Activate DHTxx sensor and collect data sent by sensor for futher processing.
static int dial_DHTxx_and_read(int32_t pin, int32_t handshakeDurUsec,
        int32_t boostPerfFlag, int32_t **arr, int32_t *arr_len, Error **err) {
    
    #if !defined(__APPLE__)
        // Set maximum priority for GPIO processing.
        if (boostPerfFlag != FALSE && -1 == set_max_priority(err)) {
            return -1;
        }
    #else
        #warning "Darwin doesn't have sched_setscheduler, so parameter boostPerfFlag is useless on Apple devices"
    #endif

    Pin p;
    if (-1 == gpio_export(pin, &p, err)) {
        gpio_unexport(&p, err);
        #if !defined(__APPLE__)
            set_default_priority(err);
        #endif            
        return -1;
    }
    // Send dial pulse.
    if (-1 == gpio_direction(&p, OUT, err)) {
        gpio_unexport(&p, err);
        #if !defined(__APPLE__)
            set_default_priority(err);
        #endif            
        return -1;
    }
    // Set pin to low.
    if (-1 == gpio_write(&p, LOW, err)) {
        gpio_unexport(&p, err);
        #if !defined(__APPLE__)
            set_default_priority(err);
        #endif            
        return -1;
    }
    // Sleep N microseconds.
    sleep_usec(handshakeDurUsec); 
    // Set pin to high.
    if (-1 == gpio_write(&p, HIGH, err)) {
        gpio_unexport(&p, err);
        #if !defined(__APPLE__)
            set_default_priority(err);
        #endif            
        return -1;
    }
    // Switch pin to input mode
    if (-1 == gpio_direction(&p, IN, err)) {
        gpio_unexport(&p, err);
        #if !defined(__APPLE__)
            set_default_priority(err);
        #endif            
        return -1;
    }
    // Read bunch of data from sensor
    // for futher processing in high level language.
    // Wait for next pulse 10ms maximum.
    if (-1 == gpio_read_seq_until_timeout(&p, 15, arr, arr_len, err)) {
        gpio_unexport(&p, err);
        #if !defined(__APPLE__)
            set_default_priority(err);
        #endif            
        return -1;
    }
    // Release pin.
    if (-1 == gpio_unexport(&p, err)) {
        #if !defined(__APPLE__)
            set_default_priority(err);
        #endif            
        return -1;
    }
    
    #if !defined(__APPLE__)
        // Return normal thread priority.
        if (boostPerfFlag != FALSE && -1 == set_default_priority(err)) {
            return -1;
        }
    #endif

    return 0;
}


// Blink specific pin n times. Led could be
// attached to this pin for debug purpose.
static int blink_n_times(int pin, int n, Error **err) {
    Pin p;
    if (-1 == gpio_export(pin, &p, err)) {
        gpio_unexport(&p, err);
        return -1;
    }
    if (-1 == gpio_direction(&p, OUT, err)) {
        gpio_unexport(&p, err);
        return -1;
    }
    int i;
    // Blink led n times in a loop.
    for (i = 0; i < n; i++)
    {
        // Turn led on
        if (-1 == gpio_write(&p, HIGH, err)) {
            gpio_unexport(&p, err);
            return -1;
        }
        // Sleep 0.1 of second.
        sleep_usec(100*1000);
        // Turn led off
        if (-1 == gpio_write(&p, LOW, err)) {
            gpio_unexport(&p, err);
            return -1;
        }
        // Sleep 0.1 of second.
        sleep_usec(100*1000);
    }
    return gpio_unexport(&p, err);
}



#endif
