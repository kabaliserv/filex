#!/bin/sh

cd web

npm install
npm run build

cp public/* dist