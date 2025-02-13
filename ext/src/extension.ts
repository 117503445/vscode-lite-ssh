// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import { get } from 'http';
import * as vscode from 'vscode';

let outputChannel: vscode.OutputChannel;

/**
 * Prints the given content on the output channel.
 *
 * @param content The content to be printed.
 * @param reveal Whether the output channel should be revealed.
 */
export const printChannelOutput = (content: string, reveal = false): void => {
	outputChannel.appendLine(content);
	if (reveal) {
		outputChannel.show(true);
	}
};

function getTerminalName(fsPath: string): string {
	return "";
	const defaultName = 'Remote Terminal';

	const splits = fsPath.split('/');
	if (splits.length <= 1) {
		return defaultName;
	}
	if (splits[0] !== 'remote') {
		return defaultName;
	}
	return splits[1];
}

// This method is called when your extension is activated
// Your extension is activated the very first time the command is executed
export function activate(context: vscode.ExtensionContext) {
	outputChannel = vscode.window.createOutputChannel("LiteSSH");

	// Use the console to output diagnostic information (console.log) and errors (console.error)
	// This line of code will only be executed once when your extension is activated
	console.log('Congratulations, your extension "litessh" is now active!');

	// The command has been defined in the package.json file
	// Now provide the implementation of the command with registerCommand
	// The commandId parameter must match the command field in package.json
	const disposable = vscode.commands.registerCommand('litessh.remoteTerminal', (resourceUri) => {
		if (!resourceUri) { return; }
		if (resourceUri.scheme === 'file') {
			const fsPath = resourceUri.fsPath;
			printChannelOutput('fsPath' + fsPath);
			const terminal = vscode.window.createTerminal(getTerminalName(fsPath));
			terminal.sendText(`r ${fsPath}`);
			terminal.show();
		}
	});

	context.subscriptions.push(disposable);
}

// This method is called when your extension is deactivated
export function deactivate() { }
