kill -9 `cat save_pid.txt` &
rm save_pid.txt &
nohup ./bin/ytsh-mail > my.log 2>&1 &
echo $! > save_pid.txt