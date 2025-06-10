#!/usr/bin/env python3

import xml.etree.ElementTree as ET
import argparse
import os

# Default file paths
report_file_path = 'junit_e2e_test.xml'
DEFAULT_INPUT = os.path.join('reports', report_file_path)
DEFAULT_OUTPUT = DEFAULT_INPUT

# Substrings to identify and filter out in testcase names
FILTER_SUBSTRINGS = [
    "[BeforeSuite]",
    "[AfterSuite]",
    "[DeferCleanup (Suite)]",
    "[DeferCleanup (Container)]",
]

# Suites to fully remove based on <testsuite name> or <testcase classname>
FILTER_SUITES = [
    "Label Selectors E2E Suite",
    "Field selectors E2E Suite",
    "Basic Operations E2E Suite",
    "Basic Operations",
    "Field Selectors Extension",        # covers partial classname matches like "Field Selectors E2E Suite"
    "Label Selectors",        # in case the classname contains it
]

def filter_junit_xml(input_path, output_path, filter_strings, suite_name_fragments):
    try:
        tree = ET.parse(input_path)
        root = tree.getroot()

        filtered_testcases_count = 0
        removed_suites_count = 0

        # Remove entire test suites by name match
        suites_to_remove = []
        for testsuite in root.findall('testsuite'):
            suite_name = testsuite.get('name', '')
            if any(keyword in suite_name for keyword in suite_name_fragments):
                suites_to_remove.append(testsuite)

        for suite in suites_to_remove:
            root.remove(suite)
            removed_suites_count += 1

        # Remove test cases based on name substrings or classname matches
        for testsuite in root.findall('testsuite'):
            testcases_to_remove = []
            for testcase in testsuite.findall('testcase'):
                name = testcase.get('name', '')
                classname = testcase.get('classname', '')
                if (any(substr in name for substr in filter_strings) or
                    any(sn in classname for sn in suite_name_fragments)):
                    testcases_to_remove.append(testcase)
                    filtered_testcases_count += 1
            for tc in testcases_to_remove:
                testsuite.remove(tc)

            # Adjust the test count in the suite
            current_tests = int(testsuite.get('tests', '0'))
            testsuite.set('tests', str(current_tests - len(testcases_to_remove)))

        tree.write(output_path, encoding='utf-8', xml_declaration=True)
        print(f"✅ Filtered {filtered_testcases_count} test case(s).")
        print(f"🗑️ Removed {removed_suites_count} entire test suite(s).")
        print(f"📄 Filtered JUnit XML saved to: {output_path}")

    except FileNotFoundError:
        print(f"❌ Error: Input file not found at {input_path}")
    except ET.ParseError as e:
        print(f"❌ Error parsing XML: {e}")
    except Exception as e:
        print(f"❌ Unexpected error: {e}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Filter unwanted test cases and test suites from a JUnit XML file.")
    parser.add_argument(
        "--input-file",
        default=DEFAULT_INPUT,
        help=f"Path to the input JUnit XML file (default: {DEFAULT_INPUT})"
    )
    parser.add_argument(
        "--output-file",
        default=DEFAULT_OUTPUT,
        help=f"Path to save the filtered output XML (default: {DEFAULT_OUTPUT})"
    )
    args = parser.parse_args()

    filter_junit_xml(args.input_file, args.output_file, FILTER_SUBSTRINGS, FILTER_SUITES)




