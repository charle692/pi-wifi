document.addEventListener('DOMContentLoaded', function() {
  getAvailableNetworks();
  document.getElementById("refresh-network-list").addEventListener('click', getAvailableNetworks);
});

function getAvailableNetworks(e) {
  if (e) e.preventDefault();

  fetch('/api/networks').then(function(response) {
    return response.json();
  }).then(function(networks) {
    var select = document.getElementById('network-selector');

    while (select.firstChild) {
      select.removeChild(select.firstChild);
    }

    networks.forEach(function(network) {
      var option = document.createElement('option');
      option.textContent = network.SSID + " - " + network.SecurityType;
      select.appendChild(option);
    });
  });
}