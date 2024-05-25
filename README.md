# malbolge

[Malbolge](https://en.wikipedia.org/wiki/Malbolge) interpreter
written in go.  Malbolge is an esoteric programming language created
by Ben Olmstead in 1998 and is known for its extreme difficulty in
writing programs due to its intentionally confusing design.

This interpreter follows the official specification from 1998 with
the "famous bug" that stops execution if the current data is outside
the 33–126 ASCII range. Malbolge Unshackled features are not
supported.

## Instructions

	go build -o malbolge .
	./malbolge <hello-world.mal>

_NOTE: If no arguments are passed then it will read of the standard
input and assume the data is valid malbolge code_

## License

MIT

