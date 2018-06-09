serverAddress="localhost:12345"
#verbosity="-v"
verbosity="--trace-ascii /dev/stdout"

function commentedOut1 {
	echo -e "\x1B[1;33;43m#####\x1B[0m"
	curl $verbosity http://$serverAddress

	echo -e "\x1B[1;33;43m#####\x1B[0m"
	curl $verbosity http://$serverAddress/

	echo -e "\x1B[1;33;43m#####\x1B[0m"
	curl $verbosity http://$serverAddress/abc
	echo -e "\x1B[1;33;43m#####\x1B[0m"
	curl $verbosity http://$serverAddress/abc/

	echo -e "\x1B[1;33;43m#####\x1B[0m"
	curl $verbosity http://$serverAddress/hash

	echo -e "\x1B[1;33;43m#####\x1B[0m"
	curl $verbosity http://$serverAddress/hash/1

		# Note about "curl -d..." and "curl --data ...":
		#	curl -d doesn't need a space after it so if you run "curl -data ..." it will be interpreted as "curl -d ata" and will send "ata" instead of the intended data
	curl $verbosity --data "password=angryMonkey" http://$serverAddress/hash
}

echo 'curl '$verbosity' --data "password=angryMonkey" --data "param2=ABC" http://'$serverAddress'/hash'
echo 'curl '$verbosity' --data "password=angryMonkey" --data "param2=ABC" http://'$serverAddress'/hash/'
echo 'curl '$verbosity' --data "password=angryMonkey&param2=ABC" http://'$serverAddress'/hash'
echo 'curl '$verbosity' --data "password=angryMonkey&param2=ABC" http://'$serverAddress'/hash/'
echo 'curl '$verbosity' "http://'$serverAddress'/hash?password=angryMonkey&param2=ABC"'
echo 'curl '$verbosity' "http://'$serverAddress'/hash/?password=angryMonkey&param2=ABC"'
echo
echo 'curl '$verbosity' "http://'$serverAddress'/hash//%2f%2F///?password=angryMonkey&param2=ABC%20DEF"'
echo 'curl '$verbosity' --data "password=angryMonkey&param2=ABC%20DEF" http://'$serverAddress'/hash//%2f%2F///'
#-X POST
