namespace Cancelamento.Demos;

/// <summary>
/// Example 2: cancellation that vanishes halfway down the call chain.
///
/// caller -> servico -> repositorio -> HttpClient. The caller gives up at
/// 300ms. In the first run every layer hands the token along and the request
/// is really interrupted. In the second run servico passes
/// CancellationToken.None "just to make it compile", and the request runs for
/// the full 2s even though nobody is waiting anymore. Nothing warns you at
/// runtime: the only clue is the server log.
/// </summary>
public static class Propagacao
{
    public static async Task RodarAsync()
    {
        await Apoio.ExigirServidorAsync();

        await RodarAsync("token repassado em todas as camadas", ServicoCorretoAsync);
        Console.WriteLine();
        await RodarAsync("servico passa CancellationToken.None", ServicoQuebradoAsync);
    }

    private static async Task RodarAsync(string titulo, Func<CancellationToken, Task<int>> servico)
    {
        Console.WriteLine($"-- {titulo} --");
        var relogio = Apoio.Cronometrar();

        using var cts = new CancellationTokenSource(TimeSpan.FromMilliseconds(300));

        try
        {
            var status = await servico(cts.Token);
            Console.WriteLine($"  resultado: status {status}");
        }
        catch (OperationCanceledException erro)
        {
            Console.WriteLine($"  resultado: {erro.GetType().Name}");
        }

        Console.WriteLine($"  IsCancellationRequested de quem chamou: {cts.IsCancellationRequested}");
        relogio.Fim("  tempo ate a chamada voltar");
    }

    private static Task<int> ServicoCorretoAsync(CancellationToken ct) => RepositorioAsync(ct);

    // The shortcut from the article. It compiles clean: nothing forces the
    // received token to be used.
    private static Task<int> ServicoQuebradoAsync(CancellationToken ct) =>
        RepositorioAsync(CancellationToken.None);

    private static async Task<int> RepositorioAsync(CancellationToken ct)
    {
        using var res = await Apoio.Http.GetAsync($"{Apoio.BaseUrl}/delay/2000", ct);
        return (int)res.StatusCode;
    }
}
