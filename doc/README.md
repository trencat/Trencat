# How to compile documentation

First change to the documentation directory

    cd doc

Build the image with `Sphinx` dependencies

    sudo docker image build -t trencat_doc:poc .

Run the container that builds the docs. Parameter `--rm` deletes the container after execution. Parameter `-v` mounts the `doc` directory of the host to the `/trencat_doc` directory of the container.

    sudo docker container run --rm -v $PWD:/trencat_doc trencat_doc:poc sphinx-build -b html /trencat_doc/source /trencat_doc/build/html

Documentation is compiled in `html` inside `./build/html`. Run the last command every time you want to build the docs again.

## Note

Documentation uses MathJax (<https://www.mathjax.org/>) to render formulas in the web browser. You can compile documentation offline, but to render the formulas in the web browser you will need internet connection. If you want to work totally offline, you can download MathJax files and follow the Sphinx instructions (https://www.sphinx-doc.org/en/master/usage/extensions/math.html).
