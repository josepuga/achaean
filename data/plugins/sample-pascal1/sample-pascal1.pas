program TcpScanPlugin;
{$MODE OBJFPC}
uses
  SysUtils, Classes, 
  Dos; // Para GetEnv

const
  PLUGIN_ID = 'sample-pascal1';

var
  progressPipe: Text;
  pluginId: String;
  progressPipeFile: String;
  paramsString: String;
  i, p: Integer;

begin
  // These values are created by Achaean
  pluginId := GetEnv('PLUGIN_ID');
  progressPipeFile := GetEnv('PROGRESS_PIPE');

  WriteLn('Hello from plugin ', pluginId);
  Flush(Output); // Flush stdout to ensure the content is written
  paramsString := '';
  for i := 1 to ParamCount do
  begin
    if i > 1 then
        paramsString := paramsString + ParamStr(i) + ' ';
  end;
  WriteLn('Parameters: ', paramsString);  
  Flush(Output);

  // Open the named pipe for progress
  Assign(progressPipe, progressPipeFile);
  try
    Rewrite(progressPipe);
  except
    on E: Exception do
    begin
      WriteLn(ErrOutput, 'Error opening named pipe: ', E.Message);
      Flush(ErrOutput); // Flush stderr
      Exit;
    end;
  end;

  // Dummy progress using while loop. Increment by 5 each iteration.
  p := 0;
  while p <= 100 do
  begin
    if p = 50 then
    begin
      WriteLn('Half of the progress...'); // Output to Stdout.
      Flush(Output); // Flush stdout

      WriteLn(ErrOutput, 'This is a fake error.'); // Output to Stderr.
      Flush(ErrOutput); // Flush stderr
    end;

    if p = 100 then
    begin
      WriteLn('Hasta la vista baby.');
      Flush(Output); // Flush stdout
    end;

    // Write progress without line breaks, using space as delimiter.
    Write(progressPipe, p, ' '); // IMPORTANT: Space required!
    Flush(progressPipe); // Ensure the value is actually written.
    Sleep(250);
    p := p + 5;
  end;

  // IMPORTANT: Plugins must end writing "DONE" to the progressPipe before exit.
  Write(progressPipe, 'DONE '); // IMPORTANT: Space required!
  Flush(progressPipe);
  Close(progressPipe);

end.
