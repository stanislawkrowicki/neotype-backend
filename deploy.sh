if [ ! "$target" ]
then
  SERVICES="web-api words users results-publisher results-consumer leaderboards"
else
  SERVICES=$target
fi

echo "Getting ready to deploy ${SERVICES}"

for SERVICE in ${SERVICES}; do
  echo "Deploying ${SERVICE}"
  heroku buildpacks:add -a "${SERVICE}-neotype" heroku-community/multi-procfile
  heroku buildpacks:add -a "${SERVICE}-neotype" heroku/go
  heroku config:set -a "${SERVICE}-neotype" PROCFILE=/cmd/"${SERVICE}"/Procfile
  git push https://git.heroku.com/"${SERVICE}-neotype".git HEAD:master
done
