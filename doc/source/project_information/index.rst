.. _project-information:

###################
Project information
###################

Getting Started
===============
This section will be filled in the future.

Prerequisites
-------------
This project is based in Python 3.6.

Installing
----------
This section will be filled in the future.

Running the tests
-----------------
Explain how to run the automated tests for this system

Compile documentation
---------------------

.. code-block:: bat

   :: Install virtualenv if you don't have it yet
   pip install virtualenv
   
   :: Create a virtualenv just for documentation libraries
   virtualenv virtualenv_doc
   
   :: Activate it
   virtualenv_doc\\Scripts\\activate.bat
   
   :: Install Sphinx and dependencies
   pip install -r doc\\requirements.txt
   
   :: Build html documentation in doc\\build\\html
   cd doc
   sphinx-build -b html source build\\html
   
.. code-block:: bash

   # Install pip if you don't have it yet
   sudo apt install python3-pip

   # Install virtualenv if you don't have it yet
   sudo pip3 install virtualenv

   # Create a virtualenv just for documentation libraries
   virtualenv virtualenv_doc

   # Activate it
   source virtualenv_doc/bin/activate

   # Install Sphinx and dependencies
   pip install -r doc/requirements.txt

   # Build html documentation in doc/build/html
   cd doc
   sphinx-build -b html source build/html


Deployment
==========
This section will be filled in the future.

.. _project-information-contributing:

Contributing
============
There are many ways you can contribute to this project. Developers, scientists, engineers, designers, geeks... Everyone can collaborate in their field of expertise. Please read `CONTRIBUTING.md <https://github.com/Joptim/Trencat/blob/master/CONTRIBUTING.md>`_ to know how can you collaborate, how can you benefit with your contributions and the details on our code of conduct.

Authors
=======
This section will be filled in the future.

License
=======
I strongly believe in high quality open source software. This software is `licensed <https://github.com/Joptim/Trencat/blob/master/LICENSE>`_ under the `GNU General Public License v3.0 <https://choosealicense.com/licenses/gpl-3.0/>`_,  which basically means that you can do almost anything you want with your project, except to distribute closed source versions.

Documentation
=============
The documentation that you are reading now introduces *TrenCAT* to interested users. The first part is devoted to give a concise state-of-the-art introduction to all subjects involved in train automation. The second part gives insights about how *TrenCAT* is structured with no specific programming language implementation details.

.. note::
	As `Carlos Ivorra <https://www.uv.es/=ivorra/>`_ says, *Mathematics not written in LaTeX are not serious mathematics*. Writing mathematics in HTML is more painful than writing in LaTeX. Seriously, just check the source of :ref:`conflict-resolution-problem-model` and you'll understand it. Docs are written in raw HTML because `solutions provided <http://www.sphinx-doc.org/es/stable/ext/math.html>`_ don't look pleasant, comfortable and lightweight enough at the time of writing.
	
Furthermore, each module has its own documentation, covering all language specific implementation details. Such implementation is generated with `Godoc <https://godoc.org/golang.org/x/tools/cmd/godoc>`_ for `Golang <https://golang.org/>`_ implementations and `Sphinx <http://www.sphinx-doc.org/en/master/>`_ for `Python <https://www.python.org/>`_ implementations.

Disclaimer
==========
The authors/sources of all mathematical theory explained in this project are explicitly mentioned in a deliberately visible position in documentation. It is important to check sources both for credits and for first hand information, better explained than it is here. If you are an author and have any concern about the content displayed in this project, please do not hesitate to contact me. Here you'll find the complete list of :ref:`citations`.

Previous topic: :ref:`Main page <main-page>`.

Next topic: :ref:`introduction-railway-infrastructure-design-theory`.