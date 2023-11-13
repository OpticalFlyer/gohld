// Create the map centered on the United States
var map = L.map('map').setView([39.8283, -98.5795], 4);

// Add the Google Maps tile layer
L.tileLayer('http://{s}.google.com/vt/lyrs=s,h&x={x}&y={y}&z={z}',{
    maxZoom: 20,
    subdomains:['mt0','mt1','mt2','mt3']
}).addTo(map);

// Initialize the draw control and pass it the FeatureGroup of editable layers
var drawnItems = new L.FeatureGroup();
map.addLayer(drawnItems);

var drawControl = new L.Control.Draw({
    draw: {
        polyline: false,
        polygon: false,
        circle: false,
        marker: false,
        circlemarker: false,
        rectangle: {
            shapeOptions: {
                color: '#f357a1',
                fillOpacity: 0.5
            }
        }
    },
    edit: false
});
map.addControl(drawControl);

// Change the tooltip of the rectangle button
var drawRectangleButton = document.querySelector('.leaflet-draw-draw-rectangle');
if (drawRectangleButton) {
    drawRectangleButton.title = 'Import road centerlines';
    drawRectangleButton.children[0].title = 'Import road centerlines';
}

var selectMode = false;

// Create a new control for the "Select Roads" button and add it to the map
var SelectControl = L.Control.extend({
    options: {
        position: 'topleft'
    },

    onAdd: function (map) {
        var container = L.DomUtil.create('div', 'leaflet-bar leaflet-control leaflet-control-custom');
        container.style.backgroundColor = 'white';
        container.style.width = '30px';
        container.style.height = '30px';
        container.style.display = 'flex';
        container.style.alignItems = 'center';
        container.style.justifyContent = 'center';
        container.style.cursor = 'pointer';

        var icon = L.DomUtil.create('i', 'icon', container);
        icon.textContent = 'S';
        icon.title = 'Select Roads';

        container.onclick = function () {
            selectMode = !selectMode;
            icon.style.color = selectMode ? 'red' : '#404040';
        };

        return container;
    }
});

map.addControl(new SelectControl());

// Create a new control for the "Delete Selected Roads" button and add it to the map
var DeleteControl = L.Control.extend({
    options: {
        position: 'topleft'
    },

    onAdd: function (map) {
        var container = L.DomUtil.create('div', 'leaflet-bar leaflet-control leaflet-control-custom');
        container.style.backgroundColor = 'white';
        container.style.width = '30px';
        container.style.height = '30px';
        container.style.display = 'flex';
        container.style.alignItems = 'center';
        container.style.justifyContent = 'center';
        container.style.cursor = 'pointer';

        var icon = L.DomUtil.create('i', 'icon', container);
        icon.textContent = 'D';
        icon.title = 'Delete Selected Roads';

        container.onclick = function () {
            map.eachLayer(function (layer) {
                if (layer instanceof L.Path && layer.options.color === 'red') {
                    map.removeLayer(layer);
                }
            });
        };

        return container;
    }
});

map.addControl(new DeleteControl());

// Create a new control for the "Download GeoJSON" button and add it to the map
var DownloadControl = L.Control.extend({
    options: {
        position: 'topleft'
    },

    onAdd: function (map) {
        var container = L.DomUtil.create('div', 'leaflet-bar leaflet-control leaflet-control-custom');
        container.style.backgroundColor = 'white';
        container.style.width = '30px';
        container.style.height = '30px';
        container.style.display = 'flex';
        container.style.alignItems = 'center';
        container.style.justifyContent = 'center';
        container.style.cursor = 'pointer';

        var icon = L.DomUtil.create('i', 'icon', container);
        icon.textContent = 'G';
        icon.title = 'Download GeoJSON';

        container.onclick = function () {
            var dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(geojson));
            var downloadAnchorNode = document.createElement('a');
            downloadAnchorNode.setAttribute("href", dataStr);
            downloadAnchorNode.setAttribute("download", "centerlines.geojson");
            document.body.appendChild(downloadAnchorNode); // required for firefox
            downloadAnchorNode.click();
            downloadAnchorNode.remove();
        };

        return container;
    }
});

map.addControl(new DownloadControl());

var geojson;

map.on('draw:created', function (e) {
var type = e.layerType,
    layer = e.layer;

if (type === 'rectangle') {
    var bounds = layer.getBounds();
    var south = bounds.getSouth();
    var west = bounds.getWest();
    var north = bounds.getNorth();
    var east = bounds.getEast();

    // Use the backend endpoint to get the road centerlines
    var requestData = {
        south: south,
        west: west,
        north: north,
        east: east
    };

    // Query the backend
    fetch('/road-centerline', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    })
    .then(response => response.json())
    .then(data => {
        geojson = data; // The backend sends back GeoJSON directly

        // Add the road centerlines to the map
        L.geoJSON(geojson, {
            onEachFeature: function (feature, layer) {
                layer.on('click', function () {
                    if (selectMode) {
                        layer.setStyle({color: 'red'});
                    } else {
                        // Generate the popup content
                        var popupContent = '<pre>' + JSON.stringify(feature.properties, null, 2) + '</pre>';

                        // Bind the popup to the layer
                        layer.bindPopup(popupContent).openPopup();
                    }
                });
            }
        }).addTo(map);

        // Remove the bounding box from the map
        map.removeLayer(layer);
    })
    .catch(err => {
        console.error('Error fetching data:', err);
    });
}

drawnItems.addLayer(layer);
});

var drawnItems = new L.FeatureGroup();
map.addLayer(drawnItems);

var drawControl = new L.Control.Draw({
    draw: {
        polyline: false,
        polygon: false, // Initially, we'll disable it here, and enable it with the custom button
        circle: false,
        marker: false,
        circlemarker: false,
        rectangle: false
    },
    edit: false
});
map.addControl(drawControl);

// Create a custom control for the "Draw Polygon to Fetch Sites" button and add it to the map
var DrawPolygonControl = L.Control.extend({
    options: {
        position: 'topleft'
    },

    onAdd: function (map) {
        var container = L.DomUtil.create('div', 'leaflet-bar leaflet-control leaflet-control-custom');
        container.style.backgroundColor = 'white';
        container.style.width = '30px';
        container.style.height = '30px';
        container.style.display = 'flex';
        container.style.alignItems = 'center';
        container.style.justifyContent = 'center';
        container.style.cursor = 'pointer';

        var icon = L.DomUtil.create('i', 'icon', container);
        icon.textContent = 'Sites';
        icon.title = 'Draw Polygon to Fetch Sites';

        container.onclick = function () {
            new L.Draw.Polygon(map, drawControl.options.draw.polygon).enable();
        };

        return container;
    }
});

map.addControl(new DrawPolygonControl());

map.on('draw:created', function (e) {
    var type = e.layerType,
        layer = e.layer;

    if (type === 'polygon') {
        var polygonCoords = layer.getLatLngs()[0];
        var requestData = {
            polygon: polygonCoords.map(coord => [coord.lat, coord.lng])
        };

        // Query the backend
        fetch('/sites-in-polygon', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestData)
        })
        .then(response => response.json())
        .then(data => {
            var sitesLayer = L.geoJSON(data, {
                pointToLayer: function(feature, latlng) {
                    return L.circleMarker(latlng, {
                        radius: 5, // size of the circle marker
                        fillColor: "#ff7800", // color of the circle marker
                        color: "#000",
                        weight: 1,
                        opacity: 1,
                        fillOpacity: 0.8
                    });
                },
                onEachFeature: function (feature, layer) {
                    if (feature.properties && feature.properties.popupContent) {
                        layer.bindPopup(feature.properties.popupContent);
                    }
                }
            }).addTo(map);

            // Remove the bounding box from the map
            map.removeLayer(layer);
        })
        .catch(err => {
            console.error('Error fetching data:', err);
        });
    }

    drawnItems.addLayer(layer);
});