#!/bin/bash
# manually exiting from script, because after-build needs to run always
set +e

tools=reference-gen
test_reporter=$PWD/cc-test-reporter

for tool in ${tools}; do
  pushd ${tool}

  if [ -z $CC_TEST_REPORTER_ID ]; then
    echo "1. CC_TEST_REPORTER_ID is unset, skipping"
  else
    echo "1. Running before-build"
    ${test_reporter} before-build
  fi

  echo "2. Running test"
  make test
  TEST_STATUS=$?

  if [ -z $CC_TEST_REPORTER_ID ]; then
    echo "3. CC_TEST_REPORTER_ID is unset, skipping"
  else
    echo "3. Running after-build"
    ${test_reporter} after-build --exit-code $TEST_STATUS -t gocov --prefix $(go list -m)
  fi

  if [ "$TEST_STATUS" -ne 0 ]; then
    echo "Test failed, status code: $TEST_STATUS"
    exit $TEST_STATUS
  fi

  popd
done
