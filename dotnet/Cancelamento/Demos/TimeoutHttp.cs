namespace Cancelamento.Demos;

/// <summary>
/// Example 1: an HTTP call canceled by a timeout.
///
/// The server takes 2s to answer and the client gives up at 300ms. Watch the
/// server log: the request shows up as "cliente desistiu", because canceling
/// the token closes the connection.
/// </summary>
public static class TimeoutHttp
{
    public static async Task RodarAsync()
    {
        await Apoio.ExigirServidorAsync();

        Console.WriteLine("-- CancellationTokenSource com timeout no construtor --");
        var relogio = Apoio.Cronometrar();

        using var cts = new CancellationTokenSource(TimeSpan.FromMilliseconds(300));

        try
        {
            using var res = await Apoio.Http.GetAsync($"{Apoio.BaseUrl}/delay/2000", cts.Token);
            Console.WriteLine($"  status: {(int)res.StatusCode}");
        }
        catch (OperationCanceledException erro)
        {
            // HttpClient throws TaskCanceledException, a subclass of
            // OperationCanceledException.
            Console.WriteLine($"  desistiu: {erro.GetType().Name}");
        }

        relogio.Fim("  tempo ate desistir");
    }
}
