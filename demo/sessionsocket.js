var sessionSocket = null;

function getAuthenticationToken() {
  return localStorage.getItem("authenticationToken");

}
function setupSessionSocket() {
  const url = 'http://0.0.0.0:8000/session';
  sessionSocket = io(url, {
    transports: ['websocket'],
  });

  sessionSocket.on('connect', () => {
    console.log('connected session socket...');

    authenticate();
  });

  sessionSocket.on('disconnect', () => {
    console.log('disconnected session socket...');
  });
}

function authenticate() {
  sessionSocket.emit('Authenticate', handleRequestData({
    roomId: 1,
  }), function (raw) {
    handleReponseData(raw, onAuthenticated, handleErrors)
  });
}

function onAuthenticated(data) {
  updatePlayer(data);

  localStorage.setItem("authenticationToken", reducer.player.authenticationToken);

  populateRooms()
}

function updatePlayer(data) {
  reducer.player = data.response.player;

  reloadUI();
}
