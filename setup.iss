[Languages]
Name: "en"; MessagesFile: "compiler:Default.isl"

[Setup]
AppName=FORUM - Cirrus Launcher
AppVerName=FORUM - Cirrus Launcher
AppContact=Carl Zeiss Meditec AG
AppCopyright=Carl Zeiss Meditec AG
AppVersion=1.0.0
AppPublisher=Carl Zeiss Meditec AG
AppPublisherURL=https://www.zeiss.de/meditec/home.html
AppendDefaultDirName=false
DefaultDirName={pf64}\CZM\CIRRUS HD-OCT\bin
DefaultGroupName=Carl Zeiss Meditec AG
VersionInfoCompany=Carl Zeiss Meditec AG
VersionInfoCopyright=Copyright © Carl Zeiss Meditec AG. All rights reserved.
VersionInfoDescription=Carl Zeiss Meditec AG
VersionInfoVersion=1.0.0
UninstallDisplayIcon={app}\cirruslauncher.exe
OutputBaseFilename=Setup CirrusLauncher
DisableReadyPage=false
AlwaysShowDirOnReadyPage=false
Compression=lzma
SolidCompression=yes
PrivilegesRequired=admin
OutputDir=.
ChangesAssociations=yes
ChangesEnvironment=yes
UsePreviousAppDir=no
UsePreviousGroup=no
DirExistsWarning=false

[Registry]
Root:HKCR; Subkey: "czmcirrus"; ValueType: string; ValueData: "URL:Custom Protocol"; Flags: uninsdeletekey; 
Root:HKCR; Subkey: "czmcirrus"; ValueType: string; ValueName: "URL Protocol"; ValueData: ""; Flags: uninsdeletekey; 
Root:HKCR; Subkey: "czmcirrus\DefaultIcon"; ValueType: string; ValueData: "{app}\CIRRUS Launcher.exe,0"; Flags: uninsdeletekey; 
Root:HKCR; Subkey: "czmcirrus\shell\open\command"; ValueType: string; ValueData: """{app}\CIRRUS Launcher.exe"" ""%1"""; Flags: uninsdeletekey; 

[Files]
Source: "D:\go\bin\cirruslauncher.exe"; DestDir: "{app}"; DestName: "CIRRUS Launcher.exe"

[UninstallDelete]
Type: files; Name: "{app}\CIRRUS Launcher.log"

[Code]
Procedure Uninstall(title : String);
var
  bool : boolean;
	cmd : String;
	ErrorCode : Integer;
begin
  bool := true;

	if IsWin64 then
	begin
		bool := not RegQueryStringValue(HKLM,'Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\' + title + '_is1','UninstallString',cmd);

		if not bool then
		begin
			bool := ShellExec('', cmd,'/VERYSILENT', '', SW_SHOWMINIMIZED, ewWaitUntilTerminated, ErrorCode);
			if not bool then
				MsgBox('Uninstall of ' + title + ' terminated with error code : ' + IntToStr(ErrorCode) + chr(10) + SysErrorMessage(ErrorCode),mbError,MB_OK);
		end;
	end
	else
	begin
		bool := not RegQueryStringValue(HKLM,'SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\' + title + '_is1','UninstallString',cmd);

		if not bool then
		begin
			bool := ShellExec('', cmd,'/VERYSILENT', '', SW_SHOWMINIMIZED, ewWaitUntilTerminated, ErrorCode);
			if not bool then
				MsgBox('Uninstall of ' + title + ' terminated with error code : ' + IntToStr(ErrorCode) + chr(10) + SysErrorMessage(ErrorCode),mbError,MB_OK);
		end;
	end;

  if not bool then
    WizardForm.Close;
end;

function NextButtonClick(PageId: Integer): Boolean;
begin
  Result := True;
  if (PageId = wpSelectDir) then 
  begin
    if not FileExists(ExpandConstant('{app}\Czm.Sdoct.AnalysisEngine.exe')) then 
    begin
      MsgBox('The selected directory does not contain the file Czm.Sdoct.AnalysisEngine.exe.  Please select the correct folder.', mbError, MB_OK);
      Result := False;
    end;
  end;                                                        
end;

procedure ChangeFile(old, nw: string);
begin
  if FileExists(old) then
  begin
    if FileExists(nw) and not DeleteFile(nw) then 
    begin
      MsgBox('Cannot delete file:' + nw, mbError, MB_OK);
      WizardForm.Close;
    end;

    if not RenameFile(old,nw) then
    begin
      MsgBox('Cannot rename file ' + old + ' to ' + nw, mbError, MB_OK);
      WizardForm.Close;
    end;
  end;
end;
