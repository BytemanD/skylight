mkdir -p /var/lib/skylight

rm -rf /usr/share/skylight /usr/local/bin/skylight
mkdir -p /usr/share /usr/local/bin || exit 1

cp -r web /usr/share/skylight  || exit 1
cp skylight /usr/local/bin || exit 1
