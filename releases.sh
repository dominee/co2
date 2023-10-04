#!/bin/bash

PROJECT_NAME="co2"

echo "Building ${PROJECT_NAME} for all platforms..."
make build-all
mv ./co2-* releases/
cd releases

echo "Compressing ${PROJECT_NAME} for all platforms..."
for BINARY_NAME in co2-darwin co2-linux co2-windows.exe
do
    zip ${BINARY_NAME}-amd64.zip ${BINARY_NAME}
done
mv co2-windows.exe-amd64.zip co2-windows-amd64.zip
zip ${PROJECT_NAME}-all-amd64.zip co2-darwin co2-linux co2-windows.exe

echo "Done!"

