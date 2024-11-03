# clear /media/uploads/1 (folder 1 may not exist)
# copy /home/my-user/core/fixtures/media/video-test/* to /media/uploads/1/*
# copy /home/my-user/core/fixtures/media/thumbnail.png to /media/uploads/thumbnail.png

# clear /media/uploads/1
rm -rf /media/uploads/*

# copy /home/my-user/core/fixtures/media/video-test/* to /media/uploads/1/*
mkdir /media/uploads/1 && cp -r /home/my-user/app/core/fixtures/media/video-test/* /media/uploads/1
mkdir /media/uploads/2 && cp -r /home/my-user/app/core/fixtures/media/video-test/* /media/uploads/2
mkdir /media/uploads/3 && cp -r /home/my-user/app/core/fixtures/media/video-test/* /media/uploads/3
mkdir /media/uploads/4 && cp -r /home/my-user/app/core/fixtures/media/video-test/* /media/uploads/4
mkdir /media/uploads/5 && cp -r /home/my-user/app/core/fixtures/media/video-test/* /media/uploads/5
mkdir /media/uploads/6 && cp -r /home/my-user/app/core/fixtures/media/video-test/* /media/uploads/6
mkdir /media/uploads/7 && cp -r /home/my-user/app/core/fixtures/media/video-test/* /media/uploads/7
mkdir /media/uploads/8 && cp -r /home/my-user/app/core/fixtures/media/video-test/* /media/uploads/8
mkdir /media/uploads/9 && cp -r /home/my-user/app/core/fixtures/media/video-test/* /media/uploads/9
mkdir /media/uploads/10 && cp -r /home/my-user/app/core/fixtures/media/video-test/* /media/uploads/10

# copy /home/my-user/core/fixtures/media/thumbnails/* to /media/uploads/* without create a folder
cp -r /home/my-user/app/core/fixtures/media/thumbnails/* /media/uploads/