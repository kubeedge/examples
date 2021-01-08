function ServiceResult() {
	var processJSONDataRaw = function (data) {
		try {
			this.Data = data;
			this.SuccessCallback.call(this);
		}
		
        catch (e) {
            // TODO: Add Err Handler. Localize Message.
            alert("Error Processing Request : " + e);
        }
	};
	
    var processJSONData = function (data) {
        try {
            this.Result = data.Result;
            this.ResultString = data.ResultString;
            this.ResultObject = data.ResultObject;

            switch (this.Result)
            {
                case 0: // Success
                    this.SuccessCallback.call(this);
                    break;

                case 100: // Validation Error
                    if (this.ValidationCallback !== undefined)
                    {
                        this.ValidationCallback.call(this);
                    }

                    break;

                case 200: // Session Timeout Error
                    alert(ResultString);
                    if (locationUrls.Logout === undefined)
                    {
                        window.location.href = '/';
                        return;
                    }

                    window.location.href = locationUrls.Logout;
                    break;

                default: // Other Error
                    alert('Error: ' + this.ResultString);
                    if (this.ErrorCallback !== undefined)
                    {
                        this.ErrorCallback.call(this);
                    }
                    break;
            }
        }
		
        catch (e) {
            // TODO: Add Err Handler. Localize Message.
            alert("Error Processing Request : " + e);
        }
    };

    // Handles Ajax Error
    var processError = function (objXHR, textStatus, error) {
        try {
            alert('Error: ' + textStatus);
            this.ErrorCallback.call(this);
        }

        catch (e) {
        }
    };

    // Object Definition
    return {
		Data: '',
        Result: '',
        ResultString: '',
        ResultObject: '',
        SuccessCallback: function () {},
        ValidationCallback: '',
        ErrorCallback: '',
        Baggage: '',

        getJSONDataBasic: function (url, data, callBack, baggage) {
            //Calls Base Method
            this.getJSONData(url, data, callback, callback, callback, baggage);
        },

        // Method to Post Data via JSON
        getJSONData: function (url, data, callback, validationCallBack, errorCallBack, baggage) {
            this.SuccessCallback = callback;
            this.ValidationCallback = validationCallBack;
            this.ErrorCallback = errorCallBack;
            this.Baggage = baggage;

            $.ajax({
                url: url,
                data: data,
                dataType: 'json',
                error: processError,
                success: processJSONData,
                context: this,
                type: 'POST'
            });
        },
		
		// Method to Post Data via JSON
        getJSONDataRaw: function (url, data, callback) {
			this.SuccessCallback = callback;
            $.ajax({
                url: url,
                data: data,
                dataType: 'json',
				error: processJSONDataRaw,
                success: processJSONDataRaw,
                context: this,
                type: 'POST'
            });
        }
    };
}