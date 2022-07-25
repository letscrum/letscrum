var common = module.exports;

common.database500 = (error) => {
  if (error) {
    // eslint-disable-next-line no-console
    console.error(error);
  }
  return {
    code: 5000,
    data: {
      message: 'Database connection error.'
    }
  };
};

common.customMessage500 = (message) => ({
  code: 5000,
  data: {
    message
  }
});

common.customData500 = (data) => ({
  code: 5000,
  data
});

common.getLibraryJoinCode = (userId, libraryId) => {
  let result = '';
  userId.split('-').forEach((s) => {
    result += s.substring(0, 1);
  });
  result += libraryId;
  return result.toUpperCase();
};

common.getLibraryByCode = (code) => {
  let result = '';
  let uIdC1 = '';
  let uIdC2 = '';
  let uIdC3 = '';
  let uIdC4 = '';
  let uIdC5 = '';
  let lId = '';
  const codeArray = [...code];
  for (let i = 0; i < codeArray.length; i += 1) {
    if (i === 0) {
      uIdC1 = codeArray[i];
    }
    if (i === 1) {
      uIdC2 = codeArray[i];
    }
    if (i === 2) {
      uIdC3 = codeArray[i];
    }
    if (i === 3) {
      uIdC4 = codeArray[i];
    }
    if (i === 4) {
      uIdC5 = codeArray[i];
    }
    if (i > 4) {
      lId += codeArray[i];
    }
  }
  result = `${uIdC1}_______-${uIdC2}___-${uIdC3}___-${uIdC4}___-${uIdC5}___________:${lId}`;
  return result.toLowerCase();
};
