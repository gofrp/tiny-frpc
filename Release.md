### Features

Introducing version 0.1.0 of our proxy solution, we enable reverse proxying through standard ssh protocol communicating with frps.

We offer two binary program options for different user needs:

The Go standalone mode: This version works independently to communicate with frps. It's built for those who favor self-sufficiency and prefer a simpler deployment process.
The ssh native mode: This version works in conjunction with the operating system's ssh program. It is intended for those whose systems comprise a preconfigured ssh setup, or who wish to utilize the existing ssh program.
Both versions necessitate the provision of a frpc toml format configuration file to function correctly.

This release is just the beginning. Stay tuned for more advanced features, improvements and we are always keen to hear user feedback for future developments.