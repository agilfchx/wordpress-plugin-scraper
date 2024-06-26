<?php

$pluginPath = $argv[1] ?? null;

if (!$pluginPath || !is_dir($pluginPath)) {
    echo "Please provide a valid plugin path.\n";
    exit(1);
}

require __DIR__ . '/vendor/autoload.php';

use PHP_CodeSniffer\Config;
use PHP_CodeSniffer\Files\FileList;
use PHP_CodeSniffer\Reporter;
use PHP_CodeSniffer\Runner;

$config = new Config(['--standard=./.phpcs.xml', $pluginPath]);

$runner = new Runner();
$runner->config = $config;
$runner->init();

$fileList = new FileList($config);
$fileList->populateFileList($config, $runner->ruleset);

if (empty($fileList->getFiles())) {
    echo "No PHP files found in the provided plugin path.\n";
    exit(0);
}

$runner->run();

$reporter = new Reporter($runner->report, $config);

$outputFile = 'reports/plugin_report_' . date('Y-m-d_H-i-s') . '.txt';
ob_start();
$reporter->printReport('full');
$reportContent = ob_get_clean();

file_put_contents($outputFile, $reportContent);
echo "Report generated: $outputFile\n";
