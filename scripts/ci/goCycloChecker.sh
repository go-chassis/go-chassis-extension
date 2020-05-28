diff -u <(echo -n) <(find . -name "*.go" -not -path "./vendor/*" -not -path ".git/*" -not -path "./third_party/*" | grep -v _test | xargs gocyclo -over 16)
if [ $? == 0 ]; then
	echo "All function has less cyclomatic complexity..."
	exit 0
else
	echo "Functions/function has more cyclomatic complexity..."
	exit 1
fi
