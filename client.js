let gameStart = 0;
let connectedPlayer = 0;

window.addEventListener("load", () => {
    const connectButton = document.getElementById("connectButton");
    const grid = document.getElementById("grid");
    let socket;

    createGrid();
    createEnemyGrid();

    // event pour le clic sur le bouton de connexion
    connectButton.addEventListener("click", () => {
        if (connectedPlayer < 2) {
            socket = initWebSocket();
            connectButton.style.display = "none";
        } else {
            console.log("There is already two player in the room");
        }
    });

    // event pour le clic sur une case de la grille
    grid.addEventListener("click", (event) => {
        if (!socket || socket.readyState !== WebSocket.OPEN) {
            console.log("Connexion WebSocket non établie !");
            return;
        }

        const clickedCell = event.target;
        const row = clickedCell.dataset.row;
        const col = clickedCell.dataset.col;

        // Envoyez les coordonnées au serveur via WebSocket
        const message = {
            type: "place_ship",
            x: parseInt(col),
            y: parseInt(row)
        };
        socket.send(JSON.stringify(message));
        socket.onmessage = (event) => {
            const response = JSON.parse(event.data);
            if (response.type === "placement_confirmation" && 
                response.message === "Bateau placé avec succès.") {
                clickedCell.classList.add("ship-placed");
                const rightCell = clickedCell.nextElementSibling;
                rightCell.classList.add("ship-placed");
                console.log("Réponse du serveur :", response.message);
            } else if (response.message === "Game start") {
                gameStart = 1;
                console.log("Game started !")
            } else {
                console.log("Bateau implacable");
            }
        };
    });
    enemyGrid.addEventListener("click", (event) => {
        if (gameStart != 1) {
            console.log("Game not started yet");
        } else {
            const clickedCell = event.target;
            const row = clickedCell.dataset.row;
            const col = clickedCell.dataset.col;

            const message = {
                type: "target_ship",
                x: parseInt(col),
                y: parseInt(row)
            };
            socket.send(JSON.stringify(message));
            socket.onmessage = (event) => {
                const response = JSON.parse(event.data);
                if (response.message === "Hit") {
                    console.log("hit")
                    clickedCell.classList.add("ship-hit");
                } else if (response.message === "Flop") {
                    console.log("flop")
                    clickedCell.classList.add("ship-missed");
                } else {
                    console.log("Error");
                }
            };
        }
    });
});

// Fonction pour créer et initialiser la connexion WebSocket
function initWebSocket() {
    socket = new WebSocket("ws://192.168.1.78:8080/ws");
    socket.onopen = () => {
        console.log("Connexion établie !");
        connectedPlayer += 1;
    };
    socket.onclose = () => {
        console.log("Connexion fermée !");
        connectedPlayer -= 1;
    };
    return socket;
}

// Fonction pour créer la grille du joueur
function createGrid() {
    for (let row = 0; row < 10; row++) {
        for (let col = 0; col < 10; col++) {
            const cell = document.createElement("div");
            cell.classList.add("cell");
            cell.dataset.row = row;
            cell.dataset.col = col;
            grid.appendChild(cell);
        }
    }
}

// Fonction pour créer la grille du joueur ennemi
function createEnemyGrid() {
    for (let row = 0; row < 10; row++) {
        for (let col = 0; col < 10; col++) {
            const cell = document.createElement("div");
            cell.classList.add("cell");
            cell.dataset.row = row;
            cell.dataset.col = col;
            enemyGrid.appendChild(cell);
        }
    }
}
