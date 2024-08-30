for((i=0;i<100;i++))
do
cast wallet new-mnemonic -a 100 >> accounts
done

grep Address accounts | awk '{print $2}' > addrs
grep "Private key" accounts | awk '{print $3}' > priks