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
  - [x] Transcode one chunk ffmpeg command
  - [x] Start the transcode for 3 chunks ahead
  - [x] If a chunk is requested that has not transcoded yet, transcode it and serve it

### Future goals

- [ ] Create the same project using DASH instead of HLS

## Getting Started

- Install ffmpeg on your machine
- Install golang on your machine
- copy (.env.example)[./.env.example] to `.env` and alter any values you like
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

## Getting started

Place a video file in your media directory (default (tmp)[./tmp]) with the name `vid.mp4`

```bash
make run

#or

make watch
```

Visit the the application on (localhost:8080)[http://localhost:8080] (unless you changed the port in .env)

In your cache directory a new (folder)[./cache/vid] will be created with a (manifest)[./cache/vid/vid.m3u8] file.
As you stream the video it will be transcoded in chunks and be played in the player

## Process

I have left this section in here as this was my documentation as I was creating the project with some of the mistakes that I made along the way with some solutions that I came up with along the way.

### Convert video to HLS format

```bash
# Without GPU
ffmpeg -i ./tmp/vid.mp4 -c:v libx264 -c:a aac -hls_list_size 0 -f hls ./public/index.m3u8

# With GPU
ffmpeg -hwaccel vaapi -vaapi_device /dev/dri/renderD128 \
  -hwaccel_output_format vaapi \
  -i ./tmp/vid.mp4 \
  -c:v h264_vaapi \
  -c:a aac \
  -hls_time 5 -hls_list_size 0 -f hls ./public/index.m3u8
```

Surprisingly you can view the video while the transcode is happening and it will reload the manifest file at the end of the video and "live stream it" which is pretty cool but it removes the ability to seek through the video and that is where you or need to wait for the entire video to be transcoded or transcode it on the fly

### Convert single chunk

```bash
ffmpeg -hwaccel vaapi -vaapi_device /dev/dri/renderD128 \
  -hwaccel_output_format vaapi \
  -ss 10 -t 10 \
  -i ./tmp/vid.mp4 \
  -c:v h264_vaapi -c:a aac \
  -f mpegts ./tmp/chunk_1.ts
```

### List all keyframes and their corresponding timestamp

With the output of the following we should be able to better predict the timestamps and precise cut points to make in the videos to ensure that each cut is on a keyframe

```bash
ffprobe -select_streams v:0 -show_frames -show_entries frame=pkt_dts_time,key_frame -of csv tmp/vid.mp4 | grep ",1"
```

The following is all of the data but just in json.

```bash
ffprobe -select_streams v:0 -show_frames -print_format json tmp/vid.mp4
```

### Current issues

At the moment I am predicting the segment durations (thinking that they will be the same as what I specified).
Turns out this is not true and I think it is due to the keyframes of a video that it cant cut at exactly the point that I have chosen.

An option that would make this entire endeavour irrelevant is to actually change that the video has more keyframes but that would require the entire video to be transcoded

An example would be to take a simple video file and generate one chunk and then do an `ffprobe` on that chunk to see what its actual duration is.
Not the number that you specified right? (If it is just know it will not always be that case for every video that you transcode)

This means that the prediction of the chunks gets all the things out of sync and it is not able to join the two chunks together.

Another issue which I can't explain yet is the fact that the first chunk when you generate all of the chunks with ffmpeg is around 10sec with the target duration of the manifest set at 10 as well even when the hls_time was set to 5.

If this is going to work I will need to find a way to correctly predict the chunk indicies and then be able to precicesly generate these chunks as the client needs them.

## Acknowledgements

- Project was created using [go-blueprint](https://github.com/Melkeydev/go-blueprint)
- For testing purposes I used [Big Buck Bunny](https://peach.blender.org/) as a test video
- [Jellyfin](https://github.com/jellyfin/jellyfin) ended up digging through some of their code to see what they ended up doing
