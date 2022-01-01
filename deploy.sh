SERVICES="web-api words"

echo "Deploying ${SERVICES}"

for SERVICE in ${SERVICES}; do
  heroku buildpacks:add -a "${SERVICE}-neotype" heroku-community/multi-procfile
  heroku buildpacks:add -a "${SERVICE}-neotype" heroku/go
  heroku config:set -a "${SERVICE}-neotype" PROCFILE=/cmd/"${SERVICE}"/Procfile
  git push https://git.heroku.com/"${SERVICE}-neotype".git HEAD:master
done
