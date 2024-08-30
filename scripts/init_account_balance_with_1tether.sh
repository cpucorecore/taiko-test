#set -x
test_account_pk=`cat main_pk`
rpc=http://192.168.100.77:28545

start_nonce=`cast nonce 0xAe95d8DA9244C37CaC0a3e16BA966a8e852Bb6D6 -r ${rpc}`
nonce=$start_nonce
echo $nonce

for to in `cat addrs`
do
cast send --async --nonce $nonce --value 1ether -r $rpc --private-key $test_account_pk $to
nonce=$((nonce+1))
done