# CHANGELOG

## v2.4.0 / 2020-01-01

- (@DABH) Node v12 support, fix node-unix-dgram issues, update dependencies
- #[115], (@pepakriz) TLS connection support
- #[123], (@JeffTomlinson, @elliotttf)  Handle oversize messages sent over UDP transports
- #[116], (@pepakriz) Make socket options configurable
- #[122], (@cjbarth) Correct improper opening and closing of sockets

## v2.2.0 / 2019-08-14

- #[82], (@AlexMost) allow use of a customer producer
- #[109], (@gdvyas) Prevent error before connection is established
- #[114], (@vrza) Support 'udp' as an alias of 'udp4'

## v2.1.0 / 2019-02-17

- (@DABH) Fix tests by fixing error emission/handling
- #[108], (@vrza) Make winston 3 a peer dependency
- #[102], (@stieg) Require winston >= 3 and add corresopnding note in readme
- #[105], (@mohd-akram) Update dependencies for latest Node compatibility

