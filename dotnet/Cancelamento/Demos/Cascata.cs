namespace Cancelamento.Demos;

/// <summary>
/// Example 4: canceling the parent cancels the child.
///
/// In Go every derived context is linked to its parent automatically. In
/// .NET the link has to be declared with CreateLinkedTokenSource. The child
/// below has its own 5s timeout, but it ends the moment the parent is
/// canceled.
/// </summary>
public static class Cascata
{
    public static async Task RodarAsync()
    {
        Console.WriteLine("-- pai cancelado, filho com timeout de 5s --");
        var relogio = Apoio.Cronometrar();

        using var pai = new CancellationTokenSource();
        using var filho = CancellationTokenSource.CreateLinkedTokenSource(pai.Token);
        filho.CancelAfter(TimeSpan.FromSeconds(5));
        using var neto = CancellationTokenSource.CreateLinkedTokenSource(filho.Token);

        using var r1 = filho.Token.Register(() => Console.WriteLine("  Register no filho"));
        using var r2 = neto.Token.Register(() => Console.WriteLine("  Register no neto"));

        await pai.CancelAsync();

        Console.WriteLine($"  filho.IsCancellationRequested: {filho.IsCancellationRequested}");
        Console.WriteLine($"  neto.IsCancellationRequested:  {neto.IsCancellationRequested}");
        relogio.Fim("  tempo (o timeout era 5s)");
    }
}
