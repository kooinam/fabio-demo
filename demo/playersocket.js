var playerSocket = null;

function setupPlayerSocket() {
  const url = '/players';
  playerSocket = io(url, {
    transports: ['websocket'],
  });

  playerSocket.on('connect', () => {
    console.log('connected player socket...');
  });

  playerSocket.on('disconnect', () => {
    console.log('disconnected player socket...');

    reducer.player = null;
  });
}

function register() {
  playerSocket.emit('Register', handleRequestData({}), function (raw) {
    handleReponseData(raw, onAuthenticated, handleError)
  });
}
