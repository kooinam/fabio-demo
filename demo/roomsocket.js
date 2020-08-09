var roomSocket = null;

function setupRoomSocket() {
  const url = 'http://0.0.0.0:8000/room';
  roomSocket = io(url, {
    transports: ['websocket'],
  });

  roomSocket.on('connect', () => {
    console.log('connected room socket...');
  });

  roomSocket.on('disconnect', () => {
    console.log('disconnected room socket...');

    reducer.room = null;
    reducer.events = [];

    reloadUI();
  });

  roomSocket.on('error', (error) => {
    console.log(error);
  });

  const eventNames = ["running", "GrabbedSeat", "LeftSeat", "MadeMove", "Waiting", "Playing", "Completed"];

  for (var i = 0; i < eventNames.length; ++i) {
    const eventName = eventNames[i];

    roomSocket.on(eventName, function (raw) {
      handleEventData(eventName, raw, function (data) {
        const event = data.event;
        reducer.events.push(event);

        updateRoom(data);

        handleRoomEvent(event);
      });
    });
  }
}

function populateRooms() {
  roomSocket.emit('List', handleRequestData({}), function (raw) {
    handleReponseData(raw, function() {

    }, handleErrors)
  });
}

function join() {
  roomSocket.emit('Join', handleRequestData({}), function (raw) {
    handleReponseData(raw, onRoomJoined, handleErrors)
  });
}

function leave() {
  roomSocket.emit('Leave', handleRequestData({}), function (raw) {
    handleReponseData(raw, onRoomLeft, handleErrors)
  });
}

function grabSeat(position) {
  roomSocket.emit('GrabSeat', handleRequestData({
    position: position,
  }), function (raw) {
    handleReponseData(raw, null, handleErrors)
  });
}

function makeMove(x, y) {
  roomSocket.emit('MakeMove', handleRequestData({
    x: x,
    y: y,
  }), function (raw) {
    handleReponseData(raw, function(data) {

    }, handleErrors)
  });
}

function send() {
  console.log(document.getElementById("message").value);
}

function updateRoom(data) {
  reducer.room = data.response.room;

  reloadUI();
}

function onRoomJoined(room) {
  updateRoom(room);
}

function onRoomLeft(room) {
  reducer.room = null;
  reducer.events = [];

  reloadUI();
}

function handleRoomEvent(event) {
  const room = reducer.room;
  const player = reducer.player;

  if (event.name === 'Playing') {
    const playerSeat = _.find(room.seats, function(seat) {
      if (seat.player && seat.player.id === player.id) {
        return true;
      }

      return false;
    });

    if (playerSeat && playerSeat.position === event.parameters.seat) {
      toastr.info('Your turn now. Make a move...');
    }
  } else if (event.name === 'Completed') {
    const winner = event.parameters.winner;

    if (winner !== -1) {
      const playerSeat = _.find(room.seats, function (seat) {
        if (seat.player && seat.player.id === player.id) {
          return true;
        }

        return false;
      });

      if (playerSeat) {
        if (winner === playerSeat.position) {
          toastr.clear();
          toastr.success('You win...');
        } else {
          toastr.clear();
          toastr.error('You lose...');
        }
      }
    } else {
      toastr.clear();
      toastr.warning('It\'s a draw...');
    }
  }
}
