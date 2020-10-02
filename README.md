# bx - Parse Burp XML

A tool to parse Burp suite HTTP proxy history XML files. 

Written in [Go](https://golang.org).

## Dependencies

- https://github.com/seh-msft/burpxml

## Build

	go build

## Usage

	; parseburpxml -h
	Usage of parseburpxml:
	-R    omit responses in CSV (as they may corrupt output in excel)
	-c    emit XML as CSV only
	-d    decode base64 bodies (may corrupt output)
	-g    emit XML as valid Go syntax only
	-i string
			input file name (rather than first argument)
	-j    emit XML as JSON only
	-o string
			output file name (rather than stdout)
	-r    omit requests in CSV (as they may corrupt output in excel)
	-s    read from stdin (rather than first argument)
	;

The file name `-` may be used to specify usage of stdin/stdout. 

## Examples

Convert XML output to a file as JSON with requests/responses decoded:

	; bx -d -j -o history.json history.xml
	; 

Get all hosts queried, omitting requests/responses from the output:

	; bx -r -R -c -i history.xml | awk -F ',' '{print $3}' | sort | uniq
	login.live.com
	outlook.office365.com
	;

Create a JSON subset of the XML data consisting of an array of paths:

	; bx -j -i history.xml | jq '.Items[] | {path: .Path}'
	{
	"path": "/async/bar"
	}
	{
	"path": "/ListStuff"
	}
	{
	"path": "/async/foo"
	}
	;
