namespace Cancelamento;

/// <summary>
/// What every example shares: the support server address, a single
/// HttpClient, a check that the server is running and a simple stopwatch.
/// </summary>
public static class Apoio
{
    /// <summary>
    /// Start the support server before the network examples:
    /// <c>node mock-server/server.js</c> or <c>go run ./cmd/00-mock-server</c>.
    /// </summary>
    public static string BaseUrl =>
        (Environment.GetEnvironmentVariable("BASE_URL") ?? "http://localhost:8080").TrimEnd('/');

    /// <summary>
    /// One HttpClient for the whole process. In a real application use
    /// IHttpClientFactory.
    /// </summary>
    public static HttpClient Http { get; } = new()
    {
        Timeout = TimeSpan.FromSeconds(30),
    };

    /// <summary>
    /// Stops the example with a clear message when the support server is not
    /// running, instead of a raw HttpRequestException. It also catches another
    /// program answering on the same port.
    /// </summary>
    public static async Task ExigirServidorAsync()
    {
        string texto;
        try
        {
            using var cts = new CancellationTokenSource(TimeSpan.FromSeconds(1));
            using var res = await Http.GetAsync($"{BaseUrl}/", cts.Token);
            texto = await res.Content.ReadAsStringAsync(cts.Token);
        }
        catch (Exception erro) when (erro is HttpRequestException or OperationCanceledException)
        {
            PararSemServidor($"O servidor de apoio nao respondeu em {BaseUrl}.");
            return;
        }

        if (!texto.StartsWith("Servidor de apoio", StringComparison.Ordinal))
        {
            PararSemServidor($"Outro programa esta respondendo em {BaseUrl}, nao o servidor de apoio.");
        }
    }

    private static void PararSemServidor(string motivo)
    {
        Console.Error.WriteLine(motivo);
        Console.Error.WriteLine("Suba o servidor de apoio antes, em outro terminal:");
        Console.Error.WriteLine();
        Console.Error.WriteLine("  node ../mock-server/server.js");
        Console.Error.WriteLine("  (ou, dentro de go/: go run ./cmd/00-mock-server)");
        Console.Error.WriteLine();
        Console.Error.WriteLine("Para usar outra porta: PORT=9000 no servidor e BASE_URL=http://localhost:9000 no exemplo.");
        Environment.Exit(1);
    }

    public static Cronometro Cronometrar() => new();
}

public sealed class Cronometro
{
    private readonly long inicio = Environment.TickCount64;

    public void Fim(string rotulo) =>
        Console.WriteLine($"{rotulo,-32} {Environment.TickCount64 - inicio}ms");
}
