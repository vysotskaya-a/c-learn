/**
 * Export source code as a .c file using client-side Blob/File API.
 */
export function exportCodeAsFile(code: string, filename = 'solution.c'): void {
  const blob = new Blob([code], { type: 'text/x-csrc;charset=utf-8' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

/**
 * Check if source code exceeds 50 KB limit.
 */
export function isCodeTooLarge(code: string): boolean {
  return new Blob([code]).size > 51200;
}
