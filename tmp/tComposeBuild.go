package tmp

const ComposeBuildTmp = `version: "3"

services:
   {{print .ModuleName}}:
      image: localhost:5000/{{print .ModuleName}}:latest
      build:
         context: ./{{print .WorkDir}}
         dockerfile: ./Dockerfile`
