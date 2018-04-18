#! /bin/bash
npm i
echo "PWD: $(pwd)"
export NODE_PATH="$NODE_PATH:$(pwd)"

run_test()
{
    TEST=$1
    PROTOFILES=$2
    shift; shift;
    [[ "$PROTOFILES" == "" ]] && PROTOFILES="tests.proto"
    echo PROTOFILES: $PROTOFILES
    docker rm -fv testtntgrpc
    (
        cd $TEST \
        && docker run -it -v $(pwd):/output -v $(pwd):/proto tntgrpc --name=tests $PROTOFILES \
        && docker run --name testtntgrpc -p3301:3301 -p50051:50051 -d \
            -v $(pwd):/opt/tarantool \
            tarantool/tarantool:2 \
            tarantool /opt/tarantool/app.lua \
        && node client.js
        CODE=$?
        docker rm -fv testtntgrpc
        if [ $CODE -eq 0 ]; then
            echo "$TEST is ok!"
        else
            echo "$TEST is failed! $CODE"
            exit $CODE
        fi
    )
}

run_test complex_object || exit 1
run_test !template_test || exit 1
run_test service_procedure_name || exit 1
run_test errors_message || exit 1
run_test errors_maps || exit 1

# for f in *; do  # iterate through diretory in "tests" path
#     if [[ -d $f && "$f" != "node_modules" ]]; then # except files and "node_modules" path
#         run_test $f
#         err=$?; [[ $err -ne 0 ]] && exit $err # exit if current test is failed
#     fi
# done
