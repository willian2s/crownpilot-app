import { frontendDirectory, run } from './process.mjs';

run('tsc', ['-b', '--pretty', 'false'], { cwd: frontendDirectory });
run('vite', ['build'], { cwd: frontendDirectory });
