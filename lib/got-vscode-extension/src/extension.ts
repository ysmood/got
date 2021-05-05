import * as vscode from 'vscode';

export function activate(context: vscode.ExtensionContext) {
	let disposable = vscode.commands.registerCommand('got-vscode-extension.testCurrent', () => {
		vscode.window.showInformationMessage('todo');
	});

	context.subscriptions.push(disposable);

	console.log('got-vscode-extension loaded');
}

export function deactivate() {}
