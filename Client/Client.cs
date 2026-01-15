using System;
using System.IO;
using System.Net.Sockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
class Client
{
    static async Task Main()
    {
        try
        {
            // Connect to server
            using var client = new TcpClient("127.0.0.1", 8080);
            using var stream = client.GetStream();
            using var reader = new StreamReader(stream, Encoding.UTF8);
            using var writer = new StreamWriter(stream, Encoding.UTF8) { AutoFlush = true };
            SendServerMessage(writer, new[] { "player", "connect" }, "||PLAYER.NAME||TestPlayer||PLAYER.DECK||CITIZEN,CITIZEN,CITIZEN||");
            // Set timeout for receiving responses
            client.Client.ReceiveTimeout = 5000;

            Console.WriteLine("Connected to server. Type 'exit' to quit.");

            // Create a cancellation token to stop the listener when user exits
            var cts = new CancellationTokenSource();

            // Start listening for server responses in the background
            var listenerTask = ListenForServerResponses(reader, cts.Token);

            // Handle user input in the foreground
            await HandleUserInput(writer, cts);

            // Wait for the listener to finish
            await listenerTask;

            Console.WriteLine("Disconnected from server.");
        }
        catch (SocketException ex)
        {
            Console.WriteLine($"Socket error: {ex.Message}");
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Error: {ex.Message}");
        }
    }

    // Continuously listen for responses from the server
    static async Task ListenForServerResponses(StreamReader reader, CancellationToken cancellationToken)
    {
        try
        {
            while (!cancellationToken.IsCancellationRequested)
            {
                var response = await reader.ReadLineAsync();
                if (response != null)
                {
                    Console.WriteLine($"\n[Server]: {response}");
                    Console.Write("\nEnter message: ");
                }
            }
        }
        catch (Exception ex)
        {
            Console.WriteLine($"[Server listener error]: {ex.Message}");
        }
    }

    // Handle user input
    static async Task HandleUserInput(StreamWriter writer, CancellationTokenSource cts)
    {
        await Task.Run(() =>
        {
            while (true)
            {
                Console.Write("\nEnter message: ");
                string? userInput = Console.ReadLine();

                if (userInput?.ToLower() == "exit")
                {
                    cts.Cancel();
                    break;
                }

                if (string.IsNullOrEmpty(userInput))
                {
                    continue;
                }

                if (userInput.ToLower() == "open lobby")
                {
                    SendServerMessage(writer, new[] { "game", "lobby", "create" }, "testLobby");
                }
                else if (userInput.ToLower() == "list lobbies")
                {
                    SendServerMessage(writer, new[] { "game", "lobby", "list" }, "");
                }
                else if (userInput.ToLower().StartsWith("join lobby "))
                {
                    string lobbyId = userInput.Substring("join lobby ".Length).Trim();
                    SendServerMessage(writer, new[] { "game", "lobby", "join" }, lobbyId);
                }
                else if (userInput.ToLower() == "start game")
                {
                    SendServerMessage(writer, new[] { "game", "lobby", "start" }, "");
                }
                else
                {
                    SendServerMessage(writer, new[] { "header1", "header2" }, userInput);
                }
            }
        });
    }

    static void SendServerMessage(StreamWriter writer, string[] headers, string message)
    {
        string header = string.Join("||HEADER.SEP||", headers) + "||HEADER.END||";
        writer.WriteLine(header + message);
    }
}
