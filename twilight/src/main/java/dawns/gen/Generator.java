package dawns.gen;

import dawns.gen.util.PropertiesFileUtil;

/**
 * 代码生成类
 */
public class Generator {

	// 根据命名规范，只修改此常量值即可
	private static String MODULE;
	private static String DATABASE;
	private static String TABLE_PREFIX;
	private static String PACKAGE_NAME;
	private static String JDBC_DRIVER;
	private static String JDBC_URL;
	private static String JDBC_USERNAME;
	private static String JDBC_PASSWORD ;
	// 需要insert后返回主键的KEY
	private static String AUTO_KEY;
	// 需要insert后返回主键的表配置，key:表名,value:主键名
//	private static Map<String, String> LAST_INSERT_ID_TABLES = new HashMap<>();
	static {
		PropertiesFileUtil propertiesFileUtil=PropertiesFileUtil.getInstance("generator");
		MODULE = propertiesFileUtil.get("generator.module");
		DATABASE = propertiesFileUtil.get("generator.jdbc.database");
		TABLE_PREFIX = propertiesFileUtil.get("generator.table.prefix");
		PACKAGE_NAME = propertiesFileUtil.get("generator.package");
		AUTO_KEY = propertiesFileUtil.get("generator.auto.key");
		JDBC_DRIVER = propertiesFileUtil.get("generator.jdbc.driver");
		JDBC_URL = propertiesFileUtil.get("generator.jdbc.url");
		JDBC_USERNAME = propertiesFileUtil.get("generator.jdbc.username");
		JDBC_PASSWORD = propertiesFileUtil.get("generator.jdbc.password");
	}

	/**
	 * 自动代码生成
	 * @param args
	 */
	public static void main(String[] args) throws Exception {
		dawns.gen.util.Generator.generator(JDBC_DRIVER, JDBC_URL, JDBC_USERNAME, JDBC_PASSWORD, MODULE, DATABASE, TABLE_PREFIX, PACKAGE_NAME, AUTO_KEY);
	}
}
