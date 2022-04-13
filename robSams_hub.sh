#sh
if [ $(pgrep robSams|wc -l) -eq 0 ]; then
  /root/go/src/robFoodDD/robSams > /var/log/robSams.log &
fi