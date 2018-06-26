serverAddress="localhost:12345"
#verbosity="-v"
verbosity="--trace-ascii /dev/stdout"
#requestType="--request POST" #"-X POST"  # <empty> = use default

rootCurlCmd="curl $verbosity $requestType"
rootURL="http://$serverAddress"

hcol="\x1B[1;33m"
rcol="\x1B[0m"

function effectivelyCommentedOut1 {
	echo -e "${hcol}#####" $rootCurlCmd $rootURL $rcol
	$rootCurlCmd                        $rootURL
	echo -e "${hcol}#####" $rootCurlCmd $rootURL/ $rcol
	$rootCurlCmd                        $rootURL/

	echo -e "${hcol}#####" $rootCurlCmd $rootURL/abc $rcol
	$rootCurlCmd                        $rootURL/abc
	echo -e "${hcol}#####" $rootCurlCmd $rootURL/abc/ $rcol
	$rootCurlCmd                        $rootURL/abc/

	echo -e "${hcol}#####" $rootCurlCmd $rootURL/hash $rcol
	$rootCurlCmd                        $rootURL/hash
	echo -e "${hcol}#####" $rootCurlCmd $rootURL/hash/ $rcol
	$rootCurlCmd                        $rootURL/hash/
	echo -e "${hcol}#####" $rootCurlCmd $rootURL/hash/1 $rcol
	$rootCurlCmd                        $rootURL/hash/1
	echo -e "${hcol}#####" $rootCurlCmd $rootURL/hash/1/ $rcol
	$rootCurlCmd                        $rootURL/hash/1/


	#################### params via --data ...
	# Note about "curl -d..." and "curl --data ...":
	#	curl -d doesn't need a space after it so if you run "curl -data ..." it will be interpreted as "curl -d ata"
	#	and will send "ata" instead of the intended data
	#
	# Note about "curl --request ...":
	#	if not specified, curl defaults to using a POST request for these

#	echo -e "${hcol}#####" $rootCurlCmd --data "param1=p1Val" $rootURL/abc $rcol
#	$rootCurlCmd                        --data "param1=p1Val" $rootURL/abc

	echo -e "${hcol}#####" $rootCurlCmd --data "param1=p1Val" --data "param2=p2Val1&param3=p3Val" --data "param2=p2Val2" $rootURL/abc $rcol
	$rootCurlCmd                        --data "param1=p1Val" --data "param2=p2Val1&param3=p3Val" --data "param2=p2Val2" $rootURL/abc
#	echo -e "${hcol}#####" $rootCurlCmd --data "param1=p1Val" --data "param2=p2Val1&param3=p3Val" --data "param2=p2Val2" $rootURL/abc/ $rcol
#	$rootCurlCmd                        --data "param1=p1Val" --data "param2=p2Val1&param3=p3Val" --data "param2=p2Val2" $rootURL/abc/


	#################### params via URL
	# Note about "curl --request ..." and Chrome & Firefox:
	#	when "--request POST" (or "-X POST") isn't specified, curl defaults to using a GET request for these which
	#	is exactly what Chrome and Firefox both send when the URL is entered into them

#	echo -e "${hcol}#####" $rootCurlCmd "$rootURL/abc?param1=p1Val" $rcol
#	$rootCurlCmd                        "$rootURL/abc?param1=p1Val"

	echo -e "${hcol}#####" $rootCurlCmd "$rootURL/abc?param1=p1Val&param2=p2Val1&param3=p3Val&param2=p2Val2" $rcol
	$rootCurlCmd                        "$rootURL/abc?param1=p1Val&param2=p2Val1&param3=p3Val&param2=p2Val2"
#	echo -e "${hcol}#####" $rootCurlCmd "$rootURL/abc/?param1=p1Val&param2=p2Val1&param3=p3Val&param2=p2Val2" $rcol
#	$rootCurlCmd                        "$rootURL/abc/?param1=p1Val&param2=p2Val1&param3=p3Val&param2=p2Val2"


	#################### params via --data ... and URL

	echo -e "${hcol}#####" $rootCurlCmd \
		--data "dataParam1=dP1Val" --data "bothParam1=bP1DVal1&dataParam2=dP2Val" --data "bothParam1=bP1DVal2" \
		"$rootURL/abc?urlParam1=uP1Val&bothParam1=bP1UVal1&urlParam2=uP2Val&bothParam1=bP1UVal2" \
		$rcol
	$rootCurlCmd \
		--data "dataParam1=dP1Val" --data "bothParam1=bP1DVal1&dataParam2=dP2Val" --data "bothParam1=bP1DVal2" \
		"$rootURL/abc?urlParam1=uP1Val&bothParam1=bP1UVal1&urlParam2=uP2Val&bothParam1=bP1UVal2"
}



# The example command from the assignment:
echo -e "${hcol}#####" $rootCurlCmd --data "password=angryMonkey" $rootURL/hash $rcol
$rootCurlCmd                        --data "password=angryMonkey" $rootURL/hash
sleep 1

echo -e "${hcol}#####" $rootCurlCmd --data "password=pass2" $rootURL/hash/ $rcol
$rootCurlCmd                        --data "password=pass2" $rootURL/hash/

echo -e "${hcol}#####" "< in a browser like Chrome, go to: $rootURL/hash?password=pass3 >" $rcol

echo -e "${hcol}#####" $rootCurlCmd $rootURL/hash/1 $rcol
$rootCurlCmd                        $rootURL/hash/1

echo -e "${hcol}#####" $rootCurlCmd $rootURL/hash/1/ $rcol
$rootCurlCmd                        $rootURL/hash/1/
sleep 5

echo -e "${hcol}#####" $rootCurlCmd $rootURL/hash/1/ $rcol
$rootCurlCmd                        $rootURL/hash/1/
sleep 1


# The example shutdown command from the assignment:
echo -e "${hcol}#####" $rootCurlCmd $rootURL/shutdown $rcol
$rootCurlCmd                        $rootURL/shutdown
	# Duplicate shutdown command to test multiple requests - this should not get a response back
	echo -e "${hcol}#####" $rootCurlCmd $rootURL/shutdown $rcol
	$rootCurlCmd                        $rootURL/shutdown


#echo
#echo 'curl '$verbosity' '$requestType' "http://'$serverAddress'/abc//%2f%2F///?password=angryMonkey&param2=ABC%20DEF"'
#echo 'curl '$verbosity' '$requestType' --data "password=angryMonkey&param2=ABC%20DEF" http://'$serverAddress'/abc//%2f%2F///'
#$rootCurlCmd "$rootURL/abc//%2f%2F///?password=angryMonkey&param2=ABC%20DEF"
