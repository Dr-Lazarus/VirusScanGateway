#!/bin/bash

# If first arg is 'bash', run bash shell
if [ "$1" = 'bash' ]; then
  exec /bin/bash
else
  # Otherwise start the server
  exec "$@"
fi
