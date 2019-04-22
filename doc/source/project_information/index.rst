.. figure:: /_static/brand/logo_horizontal.png
   :alt: TRENCAT logot.

.. _project-information:

###################
Project information
###################

Official documentation
======================
Documentation is divided into two chapters. The first part gives a concise state-of-the-art introduction to all subjects involved in train automation. The second part describes *TRENCAT*'s main structure and communications in detail.

.. Furthermore, each module has its own documentation, covering all language specific implementation details. Such implementation is generated with `Godoc <https://godoc.org/golang.org/x/tools/cmd/godoc>`_ for `Golang <https://golang.org/>`_ implementations and `Sphinx <http://www.sphinx-doc.org/en/master/>`_ for `Python <https://www.python.org/>`_ implementations.

For developers
--------------
To compile documentation, first change to the documentation directory

.. code-block:: bash

   cd doc

Build the image with `Sphinx` dependencies

.. code-block:: bash

   sudo docker image build -t trencat_doc:poc .


Run the container container that builds the docs. Parameter `--rm` deletes the container after execution. Parameter `-v` mounts the `doc` directory of the host to the `/trencat_doc` directory of the container.

.. code-block:: bash

    sudo docker container run --rm -v $PWD:/trencat_doc trencat_doc:poc sphinx-build -b html /trencat_doc/source /trencat_doc/build/html

Documentation is compiled in `html` inside `./build/html`. Run the last command every time you want to build the docs again.

.. note::

   Documentation uses `MathJax <https://www.mathjax.org/>`_ to render formulas in the web browser. You can compile documentation offline, but to render the formulas in the web browser you will need internet connection. If you want to work totally offline, you can download MathJax files and follow the Sphinx `instructions <https://www.sphinx-doc.org/en/master/usage/extensions/math.html>`_.


.. _project_information_identity_manual:

Public image guidelines
=========================
Please follow the :download:`official guidelines </_static/brand/identity_manual_en.pdf>` if you are planning to use *TRENCAT* brand publicly.

Contributing
============
There are many ways you can contribute to this project. Developers, scientists, engineers, designers, geeks... Everyone can collaborate in their field of expertise. Please read `CONTRIBUTING.md <https://github.com/Joptim/Trencat/blob/master/CONTRIBUTING.md>`_ to know how can you collaborate, how can you benefit with your contributions and the details on our code of conduct.

License
=======
I strongly believe in high quality open source software. This software is `licensed <https://github.com/Joptim/Trencat/blob/master/LICENSE>`_ under the `GNU General Public License v3.0 <https://choosealicense.com/licenses/gpl-3.0/>`_,  which basically means that you can do almost anything you want with your project, except to distribute closed source versions.

Disclaimer
==========
The authors/sources of all mathematical theory explained in this project are explicitly mentioned in a deliberately visible position in documentation. It is important to check sources both for credits and for first hand information, better explained than it is here. If you are an author and have any concern about the content displayed in this project, please do not hesitate to contact me. Here you'll find the complete list of :ref:`citations`.
