/*
 * client.js
 */

var socket = null;
var disconnected = true;

function send()
{
  var username = $('#username').val();
  var message = $('#message').val();

  var data = JSON.stringify({"user": username, "message": message});

  console.log('< sending data: ' + data);

  socket.send(data);

  $('#message').val('');  // clear sent message
  $('#message').focus();  // focus on message input
  checkInput();
}

function checkInput()
{
  // enable/disable send button
  $('#send').attr('disabled', $('#username').val().length <= 0 || $('#message').val().length <= 0 || disconnected);
}

function initClient()
{
  socket = new WebSocket('ws://localhost:8080/server')

  // on connection open
  socket.onopen = function() {
    console.log('< connected');

    var username = $('#username').val();
    socket.send('[' + username + '] connected');

    $('#messages').append(
      $('<div></div>').attr('class', 'info').text('connected to server.')
    );

    disconnected = false;
    checkInput();
  };

  // on message receive
  socket.onmessage = function(data) {
    console.log('> received data: ' + data.data);

    var received = JSON.parse(data.data);

    $('#messages').append(
      $('<div></div>').attr('class', 'message')
        .append($('<span></span>').attr('class', 'user').text(received.user))
        .append($('<span></span>').attr('class', 'message').text(received.message))
    );
  };

  // on connection close
  socket.onclose = function(e) {
    console.log('* closed with event: ' + e.code);

    $('#messages').append(
      $('<div></div>').attr('class', 'info').text('disconnected from server.')
    );

    disconnected = true;
    checkInput();
  };

  // on change input
  $(document).on('input', '#username,#message', function(){
    checkInput();
  });

  // on enter key on message input
  $('#message').keypress(function(e) {
    if(e.which == 13) // 'enter' key
    {
      $('#send').click();
      return false;  
    }
  }); 

  // initial input check
  checkInput();

  // focus on message input
  $('#message').focus();
}
