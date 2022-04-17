cd frontend
hem build
scp -r public/. brob@proj-309-yt-b-6.cs.iastate.edu:/var/www/html/NYMB
scp -r doc/. brob@proj-309-yt-b-6.cs.iastate.edu:/var/www/html/doc

ssh -tl brob proj-309-yt-b-6.cs.iastate.edu \
"cd ~/go/src/git.linux.iastate.edu/309Fall2017/YT_B_6_NYMB;\
glide install;\
go build;\
sudo mv YT_B_6_NYMB /etc/init.d;\
sudo service nymb start"
