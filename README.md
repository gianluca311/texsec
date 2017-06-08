# TexSec

TexSec is an containerized environment to safely build Latex files.

#### General
This is a PoC and work is in progress.

#### Requirements
TexSec binary depends on Docker.

To build it from source make sure you have at least Go 1.8.0 installed.

#### Installation
TexSec comes as binary or it is possible to build it from source.

To retrieve the binary, simply download the latest release for your OS and place the binary on your system.

To build it from source, run
```
go get -u github.com/gianluca311/texsec
```

##### Frontend
The frontend directory can be placed anywhere on a external webserver or even on a workstation. You only have to configure the API Host and port within the frontend/js/config.js directory.

#### Configuration
TexSec runs without any configuration. If you want to parameterized docker image, latex command/parameter and rpc point place `.texsec.yml` file into the same directory as the binary is located.
Following parameters are available with it's default values:
```
dockerImage: "tianon/latex"
daemonEndpoint: "localhost:1234"
latexCommand: "pdflatex"
latexCommandParam: ""
```

#### Team
The initial contributor team is:
* @xenolf
* @elewyyn
* @gianluca311

The initial idea for this project came from @sviehb on behalf of SEC Consult.