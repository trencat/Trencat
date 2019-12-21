Development
===========

Playbook to set Trencat development tools in Ubuntu 18.04 in localhost.


Requirements
------------

Module tested in Ubuntu 18.04.

Role Variables
--------------

Check defaults/main.yml to see a list of all settable variables and their default values.


Example Playbook
----------------

- hosts: localhost
  vars:
      # Overwrite install\_vscode variable
      - install\_vscode: true
  roles:
      - role: development

License
-------

GNU General Public License v3.0

Author Information
------------------

Check [https://github.com/trencat/Trencat](https://github.com/trencat/Trencat)
