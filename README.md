# PerfMeasurer(WIP)

Golang library to allow for easy hitting of endpoints and allowing to set headers, body, method etc

## Input

The system takes in two kinds of files, either json or yaml files can provide it by jsonfile or yamlfile flag.
Or via filename flag and specifying that the file is json or yaml using the isJson or isYaml file.
Check out the format of the json or yaml files in examples folder or in the RequestsFile object in parser folder

## Robots.txt

Comes with a Robots.txt parser  in requests folder, mostly based off of <https://github.com/samclarke/robots-parser/blob/master/Robots.js>
Allows for determination of whether or not something is allowed or not.

## Output

It will output the results based on the "-output" flag which is ./result.json by default if you supply a .yml or .yaml
file it will output the result into yaml otherwise it will make the assumption that it is outputting to json
