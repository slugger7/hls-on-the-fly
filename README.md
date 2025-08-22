# Project hls-on-the-fly

I have a unique requirement on one of my projects where I have a ton of video files that I am hosting on my application however they are all in different formats and qualities. I do not necessarily want to convert them all to HLS format and delete the originals as that would just take too long and sounds like a nightmare.

Instead what I would like to achieve is to have a way to serve up these videos that are transcoded to HLS on the fly (as the user is watching them) and afterwards keep the transcoded files for a bit and eventually just deleting them.

The service is not intended to be hosted to many people so the performance hit of the transcodes will only be local to that of the system owner.

## Goals

- [x] Set up http server
- [ ] Serve up a video file
- [ ] Serve up a video file with HLS
- [ ] Serve up a video file on the fly with HLS

### Future goals

- [ ] Create the same project using DASH instead of HLS

## Getting Started

- Install ffmpeg on your machine
- Install golang on your machine
- Use the commands in the makefile section to run the project

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

## Acknowledgements

- Project was created using [go-blueprint](https://github.com/Melkeydev/go-blueprint)