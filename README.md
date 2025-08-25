# Project hls-on-the-fly

I have a unique requirement on one of my projects where I have a ton of video files that I am hosting on my application however they are all in different formats and qualities. I do not necessarily want to convert them all to HLS format and delete the originals as that would just take too long and sounds like a nightmare.

Instead what I would like to achieve is to have a way to serve up these videos that are transcoded to HLS on the fly (as the user is watching them) and afterwards keep the transcoded files for a bit and eventually just deleting them.

The service is not intended to be hosted to many people so the performance hit of the transcodes will only be local to that of the system owner.

## Goals

- [x] Set up http server
- [x] Serve up a video file
- [x] Serve up a video file with HLS
- [ ] Serve up a video file on the fly with HLS
  - [x] Pre-create the entire .m3u8 manifest
  - [ ] Transcode one chunk ffmpeg command
  - [ ] Start the transcode for 3 chunks ahead
  - [ ] If a chunk is requested that has not transcoded yet, transcode that one and 3 chunks ahead.

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

## Convert video to HLS format

```bash
# Without GPU
ffmpeg -i ./tmp/vid.mp4 -c:v libx264 -c:a aac -hls_list_size 0 -f hls ./public/index.m3u8

# With GPU
ffmpeg -hwaccel vaapi -vaapi_device /dev/dri/renderD128 \
  -hwaccel_output_format vaapi \
  -i ./tmp/vid.mp4 \
  -c:v h264_vaapi \
  -c:a aac \
  -hls_time 10 -hls_list_size 0 -f hls ./public/index.m3u8
```

Surprisingly you can view the video while the transcode is happening and it will reload the manifest file at the end of the video and "live stream it" which is pretty cool but it removes the ability to seek through the video and that is where you or need to wait for the entire video to be transcoded or transcode it on the fly

## Acknowledgements

- Project was created using [go-blueprint](https://github.com/Melkeydev/go-blueprint)
- For testing purposes I used [Big Buck Bunny](https://peach.blender.org/) as a test video