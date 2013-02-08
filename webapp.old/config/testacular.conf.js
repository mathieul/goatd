basePath = '../';

files = [
  JASMINE,
  JASMINE_ADAPTER,
  'public/js/vendor/angular/angular.js',
  'public/js/vendor/angular-*.js',
  'test/javascripts/lib/angular/angular-mocks.js',
  'public/js/**/*.js',
  'test/javascripts/unit/**/*.js'
];

autoWatch = true;

browsers = ['Chrome'];

junitReporter = {
  outputFile: 'test_out/unit.xml',
  suite: 'unit'
};
