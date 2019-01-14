docker run -v ${PWD}:/go/src/app -w="/go/src/app" --entrypoint "/go/src/app/build.sh" maibornwolff/elcep:builder-1.10.2
Copy-Item plugin-total.so ../plugins
